package base

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	log "github.com/sirupsen/logrus"
)

const (
	FilterActionRead   = "READ"
	FilterActionCreate = "CREATE"
	FilterActionUpdate = "UPDATE"
	FilterActionDelete = "DELETE"

	OperationEqual            = "eq"
	OperationNotEqual         = "ne"
	OperationStartsWith       = "startswith"
	OperationContains         = "contains"
	OperationNotContains      = "notcontains"
	OperationExists           = "exists"
	OperationGreaterThan      = "gt"
	OperationLessThan         = "lt"
	OperationGreaterThanEqual = "ge"
	OperationLessThanEqual    = "le"
	OperationBetween          = "between"
	OperationIn               = "in"
	OperationNotIn            = "notin"
)

// OperationMap : Map of magic operator to a callable to perform the filter
type OperationMap map[string]func(string, interface{}) (interface{}, error)

type ScoutrFilters interface {
	And(key interface{}, value interface{}) interface{}
	Or(key interface{}, value interface{}) interface{}
	Equals(key string, value interface{}) (interface{}, error)
	NotEqual(key string, value interface{}) (interface{}, error)
	StartsWith(key string, value interface{}) (interface{}, error)
	Contains(key string, value interface{}) (interface{}, error)
	NotContains(key string, value interface{}) (interface{}, error)
	Exists(key string, value interface{}) (interface{}, error)
	GreaterThan(key string, value interface{}) (interface{}, error)
	LessThan(key string, value interface{}) (interface{}, error)
	GreaterThanEqual(key string, value interface{}) (interface{}, error)
	LessThanEqual(key string, value interface{}) (interface{}, error)
	Between(key string, values interface{}) (interface{}, error)
	In(key string, values interface{}) (interface{}, error)
	NotIn(key string, values interface{}) (interface{}, error)
}

// FilterBase : Interface used to generalize the filter logic across multiple providers
type FilterBase interface {
	Operations() OperationMap

	// Filter operation, Returns generated conditions and any errors
	Filter(user *types.User, filters map[string][]string, action string) (interface{}, error)

	// User filters
	userFilters(filterFields []types.FilterField) (interface{}, error)
}

type Filtering struct {
	FilterBase
	ScoutrFilters
}

func (f *Filtering) Filter(user *types.User, filters map[string][]string, action string) (interface{}, error) {
	if action == "" {
		action = FilterActionRead
	}

	var filterFields []types.FilterField
	var conditions interface{}
	var err error

	if user != nil {
		// Select filter type (defaults to read filters)
		switch action {
		case FilterActionCreate:
			filterFields = user.CreateFilters
		case FilterActionUpdate:
			filterFields = user.UpdateFilters
		case FilterActionDelete:
			filterFields = user.DeleteFilters
		default:
			filterFields = user.ReadFilters
		}

		// Perform user filter
		conditions, err = f.FilterBase.userFilters(filterFields)
		if err != nil {
			return nil, err
		}
	}

	return f.filter(conditions, filters)
}

func (f *Filtering) filter(conditions interface{}, filters map[string][]string) (interface{}, error) {
	var err error
	for key, values := range filters {
		if len(values) == 1 {
			// Perform a single query
			item := values[0]
			conditions, err = f.performFilter(conditions, key, item)
			if err != nil {
				return nil, err
			}
		} else if len(values) > 1 {
			// Perform an OR query against all possible values for this key
			// This ensures that all operations against the same key use OR operations
			var filterConds interface{}
			for _, item := range values {
				// Perform filter for this item
				result, err := f.performFilter(nil, key, item)
				if err != nil {
					return nil, err
				}

				// Combine with filterConds using OR expression
				filterConds = f.ScoutrFilters.Or(filterConds, result)
			}

			// Combine filterConds with conditions using AND expression
			conditions = f.ScoutrFilters.And(conditions, filterConds)
		}
	}

	return conditions, nil
}

func (f *Filtering) userFilters(filterFields []types.FilterField) (interface{}, error) {
	// Merge all possible values for this filter key together
	filters := make(map[string][]types.FilterField)
	for _, item := range filterFields {
		filters[item.Field] = append(filters[item.Field], item)
	}

	// Perform the filter
	var conditions interface{}
	var err error
	for key, filterItems := range filters {
		if len(filterItems) == 1 {
			// Perform a single query
			item := filterItems[0]
			conditions, err = f.performFilter(conditions, fmt.Sprintf("%s__%s", key, item.Operator), item.Value)
			if err != nil {
				return nil, err
			}
		} else if len(filterItems) > 1 {
			// Perform an OR query against all possible values for this key
			var filterConds interface{}
			for _, item := range filterItems {
				result, err := f.performFilter(nil, fmt.Sprintf("%s__%s", key, item.Operator), item.Value)
				if err != nil {
					return nil, err
				}
				filterConds = f.ScoutrFilters.Or(filterConds, result)
			}
			conditions = f.ScoutrFilters.And(conditions, filterConds)
		}
	}

	return conditions, nil
}

func (f *Filtering) getOperator(key string) (string, string) {
	// Check if this is a magic operator
	operation := OperationEqual
	matches := regexp.MustCompile("^(.+)__(.+)$").FindAllStringSubmatch(key, -1)

	if len(matches) > 0 && len(matches[0]) == 3 {
		key = matches[0][1]
		operation = matches[0][2]
	}

	return key, operation
}

func (f *Filtering) performFilter(conditions interface{}, key string, value interface{}) (interface{}, error) {
	// TODO: Unquote the value

	// Get operator
	key, operator := f.getOperator(key)

	// TODO: Convert to decimal if this is a numeric operation

	// Perform the filter operation
	fn, ok := f.FilterBase.Operations()[operator]
	if !ok {
		return nil, &types.BadRequest{
			Message: fmt.Sprintf("Provider does not support magic operator '%s'", operator),
		}
	}

	// Run the condition function
	condition, err := fn(key, value)
	if err != nil {
		return nil, err
	}

	// Apply the condition if it has a non-nil value
	if condition != nil {
		return f.ScoutrFilters.And(conditions, condition), nil
	}

	return conditions, nil
}

// MultiFilter : Build multi-filter using the IN operator to search a key for 1 or more values
func (f *Filtering) MultiFilter(user *types.User, key string, values []string) (interface{}, error) {
	// Build the default user filters
	conditions, err := f.Filter(user, nil, FilterActionRead)
	if err != nil {
		return nil, err
	}

	// Ensure the IN operation is supported by this provider
	if _, ok := f.FilterBase.Operations()[OperationIn]; !ok {
		return nil, &types.BadRequest{
			Message: "Provider does not support the IN operator",
		}
	}

	// Values are expected to be in JSON-encoded string
	vals, err := json.Marshal(values)
	if err != nil {
		log.Errorln("Failed to marshal search data")
		return nil, err
	}

	// Build the IN condition
	expr, err := f.ScoutrFilters.In(key, string(vals))
	if err != nil {
		return nil, err
	}

	conditions = f.ScoutrFilters.And(conditions, expr)

	return conditions, nil
}

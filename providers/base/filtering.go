package base

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"

	"github.com/MichaelPalmer1/scoutr-go/models"
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
type OperationMap map[string]func(string, interface{}) interface{}

// BaseFiltering : Interface used to generalize the filter logic across multiple providers
type BaseFiltering interface {
	// Implementation required by inheriting structs
	Operations() OperationMap
	And(interface{}, interface{}) interface{}
	Or(interface{}, interface{}) interface{}
	Equals(string, interface{}) interface{}

	Filter(*models.User, map[string][]string, string) (interface{}, error)
	userFilters([]models.FilterField) (interface{}, error)

	// Optional
	NotEqual(string, interface{}) interface{}
	StartsWith(string, interface{}) interface{}
	Contains(string, interface{}) interface{}
	NotContains(string, interface{}) interface{}
	Exists(string, interface{}) interface{}
	GreaterThan(string, interface{}) interface{}
	LessThan(string, interface{}) interface{}
	GreaterThanEqual(string, interface{}) interface{}
	LessThanEqual(string, interface{}) interface{}
	Between(string, interface{}) interface{}
	In(string, interface{}) interface{}
	NotIn(string, interface{}) interface{}
}

type Filtering struct {
	BaseFiltering
}

func (f *Filtering) Filter(user *models.User, filters map[string][]string, action string) (interface{}, error) {
	if action == "" {
		action = FilterActionRead
	}

	var filterFields []models.FilterField
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
		conditions, err = f.userFilters(filterFields)
		if err != nil {
			return nil, err
		}
	}

	return f.filter(conditions, filters)
}

func (f *Filtering) filter(conditions interface{}, filters map[string][]string) (interface{}, error) {
	var err error
	for key, values := range filters {
		for _, item := range values {
			// Perform the filter
			conditions, err = f.performFilter(conditions, key, item)
			if err != nil {
				return nil, err
			}
		}
	}

	return conditions, nil
}

func (f *Filtering) userFilters(filterFields []models.FilterField) (interface{}, error) {
	// Merge all possible values for this filter key together
	filters := make(map[string][]models.FilterField)
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
				filterConds = f.Or(filterConds, result)
			}
			conditions = f.And(conditions, filterConds)
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
	fn, ok := f.Operations()[operator]
	if !ok {
		return nil, &models.BadRequest{
			Message: fmt.Sprintf("Provider does not support magic operator '%s'", operator),
		}
	}

	// Run the condition function
	condition := fn(key, value)

	// If the condition is nil, do not apply the condition
	if condition != nil {
		return f.And(conditions, condition), nil
	}

	return conditions, nil
}

// Filter : Build a filter
func (api *Scoutr) OldFilter(f BaseFiltering, user *models.User, filters map[string]string) (interface{}, bool, error) {
	var conditions interface{}
	initialized := false
	re := regexp.MustCompile(`^(.+)__(.+)$`)

	// Build user filters
	if user != nil {
		for idx, item := range user.ReadFilters {
			if value, ok := item.Value.(string); ok {
				// Value is a single string
				condition := f.Equals(item.Field, value)
				initialized = true
				if idx == 0 {
					conditions = condition
				} else {
					conditions = f.And(conditions, condition)
				}
			} else if value, ok := item.Value.([]interface{}); ok {
				// Value is a list of strings
				// Check that the IN operation is supported
				if _, ok := f.Operations()["in"]; !ok {
					return nil, false, &models.BadRequest{
						Message: "Failed to generate user condition - IN operation is not supported by this provider.",
					}
				}

				// Values are expected to be in JSON-encoded string
				vals, err := json.Marshal(value)
				if err != nil {
					log.Errorln("Failed to marshal user filter data")
					return nil, false, err
				}

				// Build condition
				condition := f.In(item.Field, string(vals))
				initialized = true
				if idx == 0 {
					conditions = condition
				} else {
					conditions = f.And(conditions, condition)
				}
			} else {
				log.Warnln("Received value of unknown type", item.Value)
				log.Warnln("Type", reflect.TypeOf(item.Value))
				continue
			}
		}
	}

	// Build specified filters
	for key, value := range filters {
		var condition interface{}

		// Check for magic operator matches
		matches := re.FindAllStringSubmatch(key, -1)
		if len(matches) > 0 && len(matches[0]) == 3 {
			key = matches[0][1]
			operation := matches[0][2]

			// Find corresponding *supported* operation for this filter class
			supported := false
			for op, function := range f.Operations() {
				if operation == op {
					supported = true

					// Run the condition function
					result := function(key, value)

					// If result is nil, do not apply the condition
					if result != nil {
						condition = result
					}

					break
				}
			}

			// If filter is not found (unsupported), throw an error
			if !supported {
				return conditions, false, &models.BadRequest{
					Message: fmt.Sprintf("Unsupported magic operator '%s'", operation),
				}
			}
		} else {
			// No magic operator matches - using equals operation
			condition = f.Equals(key, value)
		}

		if !initialized {
			// Initialize conditions
			conditions = condition
			initialized = true
		} else {
			// Merge conditions together using AND
			conditions = f.And(conditions, condition)
		}
	}

	return conditions, initialized, nil
}

// MultiFilter : Build multi-filter using the IN operator to search a key for 1 or more values
func (api *Scoutr) MultiFilter(f BaseFiltering, user *models.User, key string, values []string) (interface{}, error) {
	// Build the default user filters
	conditions, err := api.Filtering.Filter(user, nil, "")
	if err != nil {
		return nil, err
	}

	// Ensure the IN operation is supported by this provider
	if _, ok := f.Operations()["in"]; !ok {
		return nil, &models.BadRequest{
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
	conditions = f.And(conditions, f.In(key, string(vals)))

	return conditions, nil
}

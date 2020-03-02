package base

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"

	"github.com/MichaelPalmer1/simple-api-go/models"
	log "github.com/sirupsen/logrus"
)

// OperationMap : Map of magic operator to a callable to perform the filter
type OperationMap map[string]func(string, interface{}) interface{}

// Filtering : Interface used to generalize the filter logic across multiple providers
type Filtering interface {
	// Implementation required by inheriting structs
	Operations() OperationMap
	And(interface{}, interface{}) interface{}
	Equals(string, interface{}) interface{}

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

// Filter : Build a filter
func (api *SimpleAPI) Filter(f Filtering, user *models.User, filters map[string]string) (interface{}, bool, error) {
	var conditions interface{}
	initialized := false
	re := regexp.MustCompile(`^(.+)__(in|notin|contains|notcontains|startswith|ne|gt|lt|ge|le|between|exists)$`)

	// Build user filters
	if user != nil {
		for idx, item := range user.FilterFields {
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
func (api *SimpleAPI) MultiFilter(f Filtering, user *models.User, key string, values []string) (interface{}, error) {
	// Build the default user filters
	conditions, hasValues, err := api.Filter(f, user, nil)
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
	condition := f.In(key, string(vals))

	if hasValues {
		// Merge conditions using an AND operation
		conditions = f.And(conditions, condition)
	} else {
		// Initialize condition
		conditions = condition
	}

	return conditions, nil
}

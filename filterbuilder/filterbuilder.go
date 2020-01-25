package filterbuilder

import (
	"encoding/json"
	"reflect"
	"regexp"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/sirupsen/logrus"
)

// Filter : Build a filter
func Filter(user *models.User, filters map[string]string) (expression.ConditionBuilder, bool) {
	var conditions expression.ConditionBuilder
	initialized := false
	re := regexp.MustCompile(`^(.+)__(in|contains|notcontains|startswith|ne|gt|lt|ge|le|between|exists)$`)

	// Build user filters
	if user != nil {
		for idx, item := range user.FilterFields {
			attr := expression.Name(item.Field)
			if value, ok := item.Value.(string); ok {
				// Value is a single string
				condition := attr.Equal(expression.Value(value))
				initialized = true
				if idx == 0 {
					conditions = condition
				} else {
					conditions = conditions.And(condition)
				}
			} else if value, ok := item.Value.([]interface{}); ok {
				// Value is a list of strings
				condition := attr.In(expression.Value(value))
				initialized = true
				if idx == 0 {
					conditions = condition
				} else {
					conditions = conditions.And(condition)
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
		condition := expression.ConditionBuilder{}

		// Check for magic operator matches
		matches := re.FindAllStringSubmatch(key, -1)
		if len(matches) > 0 && len(matches[0]) == 3 {
			key = matches[0][1]
			operation := matches[0][2]
			attr := expression.Name(key)

			// Perform filter based on the desired magic operator
			switch operation {
			case "in":
				var valueList []string
				json.Unmarshal([]byte(value), &valueList)
				condition = attr.In(expression.Value(valueList))
			// case "notin":
			// 	var valueList []string
			// 	json.Unmarshal([]byte(value), &valueList)
			// 	condition = expression.Not(attr.In(expression.Value(valueList)))
			case "contains":
				condition = attr.Contains(value)
			case "notcontains":
				condition = expression.Not(attr.Contains(value))
			case "exists":
				if value == "true" {
					condition = attr.AttributeExists()
				} else if value == "false" {
					condition = attr.AttributeNotExists()
				} else {
					continue
				}
			case "startswith":
				condition = attr.BeginsWith(value)
			case "ne":
				condition = attr.NotEqual(expression.Value(value))
			case "between":
				var valueList []string
				json.Unmarshal([]byte(value), &valueList)
				condition = attr.Between(expression.Value(valueList[0]), expression.Value(valueList[1]))
			case "gt":
				condition = attr.GreaterThan(expression.Value(value))
			case "lt":
				condition = attr.LessThan(expression.Value(value))
			case "ge":
				condition = attr.GreaterThanEqual(expression.Value(value))
			case "le":
				condition = attr.LessThanEqual(expression.Value(value))
			default:
				panic("Unsupported magic operator")
			}
		} else {
			condition = expression.Name(key).Equal(expression.Value(value))
		}

		if initialized {
			conditions = conditions.And(condition)
		} else {
			conditions = condition
			initialized = true
		}
	}

	return conditions, initialized
}

// MultiFilter : Perform a filter with multiple values
func MultiFilter(user *models.User, key string, values []string) expression.ConditionBuilder {
	conditions, hasValues := Filter(user, nil)

	// Build the condition
	condition := expression.Name(key).In(expression.Value(values))

	if hasValues {
		conditions = conditions.And(condition)
	} else {
		conditions = condition
	}

	return conditions
}

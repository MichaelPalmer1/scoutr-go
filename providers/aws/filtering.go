package aws

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/MichaelPalmer1/scoutr-go/providers/base"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

type DynamoFiltering struct {
	*base.Filtering
}

// Operations : Map of supported operations for this filter provider
func (f *DynamoFiltering) Operations() base.OperationMap {
	return base.OperationMap{
		base.OperationStartsWith:       f.StartsWith,
		base.OperationEqual:            f.Equals,
		base.OperationNotEqual:         f.NotEqual,
		base.OperationContains:         f.Contains,
		base.OperationNotContains:      f.NotContains,
		base.OperationExists:           f.Exists,
		base.OperationGreaterThan:      f.GreaterThan,
		base.OperationLessThan:         f.LessThan,
		base.OperationGreaterThanEqual: f.GreaterThanEqual,
		base.OperationLessThanEqual:    f.LessThanEqual,
		base.OperationBetween:          f.Between,
		base.OperationIn:               f.In,
		base.OperationNotIn:            f.NotIn,
	}
}

// And : Takes two conditions and performs an AND operation on them
func (f *DynamoFiltering) And(condition1, condition2 interface{}) interface{} {
	if condition1 != nil && condition2 != nil {
		return expression.And(condition1.(expression.ConditionBuilder), condition2.(expression.ConditionBuilder))
	} else if condition1 != nil {
		return condition1
	} else if condition2 != nil {
		return condition2
	} else {
		return nil
	}
}

// Or : Takes two conditions and performs an OR operation on them
func (f *DynamoFiltering) Or(condition1, condition2 interface{}) interface{} {
	if condition1 != nil && condition2 != nil {
		return expression.Or(condition1.(expression.ConditionBuilder), condition2.(expression.ConditionBuilder))
	} else if condition1 != nil {
		return condition1
	} else if condition2 != nil {
		return condition2
	} else {
		return nil
	}
}

// Equals : Standard equals operation
func (f *DynamoFiltering) Equals(key string, value interface{}) interface{} {
	return expression.Name(key).Equal(expression.Value(value))
}

// NotEqual : Standard not equals operation
func (f *DynamoFiltering) NotEqual(key string, value interface{}) interface{} {
	return expression.Name(key).NotEqual(expression.Value(value))
}

// StartsWith : Find all records that contain items that start with a specific value
func (f *DynamoFiltering) StartsWith(key string, value interface{}) interface{} {
	return expression.Name(key).BeginsWith(value.(string))
}

// Contains : Check if a value contains a string
func (f *DynamoFiltering) Contains(key string, value interface{}) interface{} {
	return expression.Name(key).Contains(value.(string))
}

// NotContains : Check for values that do not contain a string
func (f *DynamoFiltering) NotContains(key string, value interface{}) interface{} {
	return expression.Not(expression.Name(key).Contains(value.(string)))
}

// Exists : Checks if an attribute exists. Only accepts true/false values. Returns nil for all other values.
func (f *DynamoFiltering) Exists(key string, value interface{}) interface{} {
	attr := expression.Name(key)
	if value == "true" {
		return attr.AttributeExists()
	} else if value == "false" {
		return attr.AttributeNotExists()
	} else {
		return nil
	}
}

// GreaterThan : Check if a value is greater than a string
func (f *DynamoFiltering) GreaterThan(key string, value interface{}) interface{} {
	return expression.Name(key).GreaterThan(expression.Value(value))
}

// LessThan : Check if a value is greater than a string
func (f *DynamoFiltering) LessThan(key string, value interface{}) interface{} {
	return expression.Name(key).LessThan(expression.Value(value))
}

// GreaterThanEqual : Check if a value is greater than a string
func (f *DynamoFiltering) GreaterThanEqual(key string, value interface{}) interface{} {
	return expression.Name(key).GreaterThanEqual(expression.Value(value))
}

// LessThanEqual : Check if a value is greater than a string
func (f *DynamoFiltering) LessThanEqual(key string, value interface{}) interface{} {
	return expression.Name(key).LessThanEqual(expression.Value(value))
}

// Between : Check for records that are between a low and high value
func (f *DynamoFiltering) Between(key string, value interface{}) interface{} {
	var valueList []string
	err := json.Unmarshal([]byte(value.(string)), &valueList)
	if err != nil {
		log.Errorf("Failed to unmarshal data: %v", err)
		return nil
	}
	return expression.Name(key).Between(expression.Value(valueList[0]), expression.Value(valueList[1]))
}

// In : Find all records with a list of values
func (f *DynamoFiltering) In(key string, values interface{}) interface{} {
	var valueList []string
	err := json.Unmarshal([]byte(values.(string)), &valueList)
	if err != nil {
		log.Errorf("Failed to unmarshal data: %v", err)
		return nil
	}

	// Generate the IN filter for this condition
	inFilter := generateInFilter(key, valueList)

	// Make sure an expression was generated
	if inFilter != nil {
		return *inFilter
	} else {
		return nil
	}
}

// NotIn : Find all records without a list of values
func (f *DynamoFiltering) NotIn(key string, values interface{}) interface{} {
	conditions := f.In(key, values)
	if conditions != nil {
		return conditions.(expression.ConditionBuilder).Not()
	} else {
		return nil
	}
}

// generateInFilter : Given a column and list of values, generate a filter expression for checking if
// the column contains the values. Returns nil if no values are supplied.
func generateInFilter(key string, values []string) *expression.ConditionBuilder {
	var firstValue expression.OperandBuilder
	var filterValues []expression.OperandBuilder
	var condition expression.ConditionBuilder

	// nil check on provided values
	if values == nil {
		return nil
	}

	// Loop through all the values
	for n, value := range values {
		// Save the first value
		if n == 0 {
			firstValue = expression.Value(value)
			continue
		}

		// Create expression filterValues at indexes 1 to N
		filterValues = append(filterValues, expression.Value(value))
	}

	if len(values) == 0 {
		return nil
	} else if len(values) == 1 {
		condition = expression.Name(key).In(firstValue)
	} else {
		condition = expression.Name(key).In(firstValue, filterValues...)
	}

	return &condition
}

func (f *DynamoFiltering) BuildInExpr(attr string, values []string, negate bool) interface{} {
	startIndex := 0
	endIndex := 0
	var conditions interface{}

	// IN expressions are limited to 100 items each
	for i := 0; i < len(values); i += 100 {
		var expr expression.ConditionBuilder

		// Create a slice of 100 items
		items := values[startIndex:endIndex]

		// Skip if no items are in this slice
		if len(items) == 0 {
			continue
		}

		// Create IN expression
		inExpr := generateInFilter(attr, items)
		if inExpr == nil {
			startIndex = endIndex
			continue
		}
		if negate {
			expr = inExpr.Not()
		} else {
			expr = *inExpr
		}

		// Combine with conditions using OR
		if negate {
			conditions = f.And(conditions, expr)
		} else {
			conditions = f.Or(conditions, expr)
		}

		// Set new start index
		startIndex = endIndex

	}

	// Add any extra items at the end
	if len(values[endIndex:]) > 0 {
		expr := generateInFilter(attr, values[endIndex:])
		if expr == nil {
			return conditions
		}

		if negate {
			conditions = f.And(conditions, expr.Not())
		} else {
			conditions = f.Or(conditions, expr)
		}
	}

	return conditions
}

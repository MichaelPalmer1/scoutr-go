package aws

import (
	"encoding/json"
	"fmt"

	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

type DynamoFiltering struct {
	base.Filtering
}

func NewFilter() DynamoFiltering {
	f := DynamoFiltering{}
	f.FilterBase = &f
	f.ScoutrFilters = &f
	return f
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
	cond1, ok1 := condition1.(expression.ConditionBuilder)
	cond2, ok2 := condition2.(expression.ConditionBuilder)

	if ok1 && cond1.IsSet() && ok2 && cond2.IsSet() {
		return expression.And(cond1, cond2)
	} else if ok1 && cond1.IsSet() {
		return cond1
	} else if ok2 && cond2.IsSet() {
		return cond2
	} else {
		return expression.ConditionBuilder{}
	}
}

// Or : Takes two conditions and performs an OR operation on them
func (f *DynamoFiltering) Or(condition1, condition2 interface{}) interface{} {
	cond1, ok1 := condition1.(expression.ConditionBuilder)
	cond2, ok2 := condition2.(expression.ConditionBuilder)

	if ok1 && cond1.IsSet() && ok2 && cond2.IsSet() {
		return expression.Or(cond1, cond2)
	} else if ok1 && cond1.IsSet() {
		return cond1
	} else if ok2 && cond2.IsSet() {
		return cond2
	} else {
		return expression.ConditionBuilder{}
	}
}

// Equals : Standard equals operation
func (f *DynamoFiltering) Equals(key string, value interface{}) (interface{}, error) {
	return expression.Name(key).Equal(expression.Value(value)), nil
}

// NotEqual : Standard not equals operation
func (f *DynamoFiltering) NotEqual(key string, value interface{}) (interface{}, error) {
	return expression.Name(key).NotEqual(expression.Value(value)), nil
}

// StartsWith : Find all records that contain items that start with a specific value
func (f *DynamoFiltering) StartsWith(key string, value interface{}) (interface{}, error) {
	return expression.Name(key).BeginsWith(value.(string)), nil
}

// Contains : Check if a value contains a string
func (f *DynamoFiltering) Contains(key string, value interface{}) (interface{}, error) {
	return expression.Name(key).Contains(value.(string)), nil
}

// NotContains : Check for values that do not contain a string
func (f *DynamoFiltering) NotContains(key string, value interface{}) (interface{}, error) {
	return expression.Not(expression.Name(key).Contains(value.(string))), nil
}

// Exists : Checks if an attribute exists. Only accepts true/false values. Returns nil for all other values.
func (f *DynamoFiltering) Exists(key string, value interface{}) (interface{}, error) {
	attr := expression.Name(key)
	if value == "true" {
		return attr.AttributeExists(), nil
	} else if value == "false" {
		return attr.AttributeNotExists(), nil
	} else {
		return nil, fmt.Errorf("invalid value for Exists operation. Supported values are ['true'/'false']")
	}
}

// GreaterThan : Check if a value is greater than a string
func (f *DynamoFiltering) GreaterThan(key string, value interface{}) (interface{}, error) {
	return expression.Name(key).GreaterThan(expression.Value(value)), nil
}

// LessThan : Check if a value is greater than a string
func (f *DynamoFiltering) LessThan(key string, value interface{}) (interface{}, error) {
	return expression.Name(key).LessThan(expression.Value(value)), nil
}

// GreaterThanEqual : Check if a value is greater than a string
func (f *DynamoFiltering) GreaterThanEqual(key string, value interface{}) (interface{}, error) {
	return expression.Name(key).GreaterThanEqual(expression.Value(value)), nil
}

// LessThanEqual : Check if a value is greater than a string
func (f *DynamoFiltering) LessThanEqual(key string, value interface{}) (interface{}, error) {
	return expression.Name(key).LessThanEqual(expression.Value(value)), nil
}

// Between : Check for records that are between a low and high value
//
// Operator: key__between=["1", "2"]
func (f *DynamoFiltering) Between(key string, values interface{}) (interface{}, error) {
	var valueList []string

	// Convert to string
	s, ok := values.(string)
	if !ok {
		log.Errorf("Failed to cast %+v to string", values)
		return nil, fmt.Errorf("%+v could not be cast as a string", values)
	}

	// Unmarshal JSON
	err := json.Unmarshal([]byte(s), &valueList)
	if err != nil {
		log.WithError(err).Error("Failed to unmarshal data")
		return nil, err
	}

	return expression.Name(key).Between(expression.Value(valueList[0]), expression.Value(valueList[1])), nil
}

// In : Find all records with a list of values
func (f *DynamoFiltering) In(key string, values interface{}) (interface{}, error) {
	var valueList []string

	// Convert to string
	s, ok := values.(string)
	if !ok {
		log.Errorf("Failed to cast %+v to string", values)
		return nil, fmt.Errorf("%+v could not be cast as a string", values)
	}

	err := json.Unmarshal([]byte(s), &valueList)
	if err != nil {
		log.Errorf("Failed to unmarshal data: %v", err)
		return nil, err
	}

	// Generate the IN filter for this condition
	return f.BuildInExpr(key, valueList, false), nil
}

// NotIn : Find all records without a list of values
func (f *DynamoFiltering) NotIn(key string, values interface{}) (interface{}, error) {
	var valueList []string

	// Convert to string
	s, ok := values.(string)
	if !ok {
		log.Errorf("Failed to cast %+v to string", values)
		return nil, fmt.Errorf("%+v could not be cast as a string", values)
	}

	err := json.Unmarshal([]byte(s), &valueList)
	if err != nil {
		log.Errorf("Failed to unmarshal data: %v", err)
		return nil, err
	}

	// Generate the NOT IN filter for this condition
	return f.BuildInExpr(key, valueList, true), nil
}

// generateInFilter : Given a column and list of values, generate a filter expression for checking if
// the column contains the values. Returns nil if no values are supplied.
func generateInFilter(key string, values []string) expression.ConditionBuilder {
	var firstValue expression.OperandBuilder
	var filterValues []expression.OperandBuilder
	var condition expression.ConditionBuilder

	// nil check on provided values
	if values == nil {
		return condition
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
		return condition
	} else if len(values) == 1 {
		condition = expression.Name(key).In(firstValue)
	} else {
		condition = expression.Name(key).In(firstValue, filterValues...)
	}

	return condition
}

func (f *DynamoFiltering) BuildInExpr(attr string, values []string, negate bool) expression.ConditionBuilder {
	var conditions interface{}
	startIndex := 0
	endIndex := 0

	// IN expressions are limited to 100 items each
	for i := 100; i < len(values); i += 100 {
		var expr expression.ConditionBuilder
		endIndex = i

		// Create a slice of 100 items
		items := values[startIndex:endIndex]

		// Create IN expression
		inExpr := generateInFilter(attr, items)
		if !inExpr.IsSet() {
			startIndex = endIndex
			continue
		}

		// Negation
		if negate {
			expr = inExpr.Not()
		} else {
			expr = inExpr
		}

		// Combine with conditions using AND/OR
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

		if negate {
			conditions = f.And(conditions, expr.Not())
		} else {
			conditions = f.Or(conditions, expr)
		}
	}

	if conds, ok := conditions.(expression.ConditionBuilder); !ok {
		return expression.ConditionBuilder{}
	} else {
		return conds
	}
}

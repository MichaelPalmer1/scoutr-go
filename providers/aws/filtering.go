package aws

import (
	"encoding/json"

	"github.com/MichaelPalmer1/simple-api-go/providers/base"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type DynamoFiltering struct {
	base.Filtering
}

func (f *DynamoFiltering) Operations() base.OperationMap {
	return base.OperationMap{
		"startswith":  f.StartsWith,
		"ne":          f.NotEquals,
		"contains":    f.Contains,
		"notcontains": f.NotContains,
		"exists":      f.Exists,
		"gt":          f.GreaterThan,
		"lt":          f.LessThan,
		"ge":          f.GreaterThanEqual,
		"le":          f.LessThanEqual,
		"between":     f.Between,
		// "in":          f.In,
	}
}

// And : Takes two conditions and performs an AND operation on them
func (f *DynamoFiltering) And(conditions, condition interface{}) interface{} {
	return expression.And(conditions.(expression.ConditionBuilder), condition.(expression.ConditionBuilder))
}

// StartsWith : Find all records that contain items that start with a specific value
func (f *DynamoFiltering) StartsWith(key string, value interface{}) interface{} {
	return expression.Name(key).BeginsWith(value.(string))
}

// Equals : Standard equals operation
func (f *DynamoFiltering) Equals(key string, value interface{}) interface{} {
	return expression.Name(key).Equal(expression.Value(value))
}

// NotEquals : Standard not equals operation
func (f *DynamoFiltering) NotEquals(key string, value interface{}) interface{} {
	return expression.Name(key).NotEqual(expression.Value(value))
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
	json.Unmarshal([]byte(value.(string)), &valueList)
	return expression.Name(key).Between(expression.Value(valueList[0]), expression.Value(valueList[1]))
}

// In : Find all records with a list of values
func (f *DynamoFiltering) In(key string, values interface{}) interface{} {
	// TODO: Some reason, this operator does not seem to work right in Go...
	var valueList []string
	json.Unmarshal([]byte(values.(string)), &valueList)
	return expression.Name(key).In(expression.Value(valueList))
}

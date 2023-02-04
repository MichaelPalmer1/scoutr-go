package mongo

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/MichaelPalmer1/scoutr-go/providers/base"
	"github.com/globalsign/mgo/bson"
)

// MongoDBFiltering : Used by the MongoAPI to perform filters against a MongoDB backend
type MongoDBFiltering struct {
	base.FilterBase
}

// Operations : Map of supported operations for this filter provider
func (f *MongoDBFiltering) Operations() base.OperationMap {
	return base.OperationMap{
		"ne":         f.NotEqual,
		"startswith": f.StartsWith,
		"contains":   f.Contains,
		// TODO: NotContains does not work right
		//"notcontains": f.NotContains,
		"exists":  f.Exists,
		"gt":      f.GreaterThan,
		"lt":      f.LessThan,
		"ge":      f.GreaterThanEqual,
		"le":      f.LessThanEqual,
		"between": f.Between,
		"in":      f.In,
		"notin":   f.NotIn,
	}
}

// And : Takes two conditions and performs an AND operation on them
func (f *MongoDBFiltering) And(conditions, condition interface{}) interface{} {
	var output bson.D

	if _, ok := conditions.(bson.DocElem); ok {
		// If element, add the element
		output = append(output, conditions.(bson.DocElem))
	} else {
		// If element array, add all items of the array
		output = append(output, conditions.(bson.D)...)
	}

	if _, ok := condition.(bson.DocElem); ok {
		// If element, add the element
		output = append(output, condition.(bson.DocElem))
	} else {
		// If element array, add all items of the array
		output = append(output, condition.(bson.D)...)
	}

	return output
}

// Equals : Standard equals operation
func (f *MongoDBFiltering) Equals(key string, value interface{}) interface{} {
	return bson.DocElem{Name: key, Value: value}
}

// NotEqual : Standard not equals operation
func (f *MongoDBFiltering) NotEqual(key string, value interface{}) interface{} {
	return bson.DocElem{
		Name: key,
		Value: bson.D{{
			Name:  "$ne",
			Value: value,
		}},
	}
}

// StartsWith : Check if string starts with a value
func (f *MongoDBFiltering) StartsWith(key string, value interface{}) interface{} {
	return bson.DocElem{
		Name: key,
		Value: bson.D{{
			Name:  "$regex",
			Value: fmt.Sprintf("^%s.*", value),
		}},
	}
}

// Contains : Check if string is contained in record
func (f *MongoDBFiltering) Contains(key string, value interface{}) interface{} {
	return bson.DocElem{
		Name: key,
		Value: bson.D{{
			Name:  "$regex",
			Value: value,
		}},
	}
}

// TODO: Feel like this *should* work, but it never returns results
// NotContains : Check if string is not contained in record
func (f *MongoDBFiltering) NotContains(key string, value interface{}) interface{} {
	return bson.DocElem{
		Name: key,
		Value: bson.DocElem{
			Name: "$not",
			Value: bson.D{{
				Name:  "$regex",
				Value: value,
			}},
		},
	}
}

// Exists : Check if attribute exists
func (f *MongoDBFiltering) Exists(key string, value interface{}) interface{} {
	var exists bool
	if value == "true" {
		exists = true
	} else if value == "false" {
		exists = false
	} else {
		log.Warnf("Invalid value for EXISTS operation: %s", value)
		return nil
	}
	return bson.DocElem{
		Name: key,
		Value: bson.D{{
			Name:  "$exists",
			Value: exists,
		}},
	}
}

// GreaterThan : Check if a value is greater than a string
func (f *MongoDBFiltering) GreaterThan(key string, value interface{}) interface{} {
	return bson.DocElem{
		Name: key,
		Value: bson.D{{
			Name:  "$gt",
			Value: value,
		}},
	}
}

// LessThan : Check if a value is greater than a string
func (f *MongoDBFiltering) LessThan(key string, value interface{}) interface{} {
	return bson.DocElem{
		Name: key,
		Value: bson.D{{
			Name:  "$lt",
			Value: value,
		}},
	}
}

// GreaterThanEqual : Check if a value is greater than a string
func (f *MongoDBFiltering) GreaterThanEqual(key string, value interface{}) interface{} {
	return bson.DocElem{
		Name: key,
		Value: bson.D{{
			Name:  "$gte",
			Value: value,
		}},
	}
}

// LessThanEqual : Check if a value is greater than a string
func (f *MongoDBFiltering) LessThanEqual(key string, value interface{}) interface{} {
	return bson.DocElem{
		Name: key,
		Value: bson.D{{
			Name:  "$lte",
			Value: value,
		}},
	}
}

// Between : Check for records that are between a low and high value
func (f *MongoDBFiltering) Between(key string, value interface{}) interface{} {
	var valueList []string
	err := json.Unmarshal([]byte(value.(string)), &valueList)
	if err != nil {
		log.Errorf("Failed to unmarshal value list for BETWEEN operation: %v", err)
		return nil
	}
	return f.And(
		f.GreaterThanEqual(key, valueList[0]),
		f.LessThanEqual(key, valueList[1]),
	)
}

// In : Find all records with a list of values
func (f *MongoDBFiltering) In(key string, values interface{}) interface{} {
	var valueList []string
	err := json.Unmarshal([]byte(values.(string)), &valueList)
	if err != nil {
		log.Errorf("Failed to unmarshal value list for IN operation: %v", err)
		return nil
	}
	return bson.DocElem{
		Name: key,
		Value: bson.D{
			{Name: "$in", Value: valueList},
		},
	}
}

// NotIn : Find all records not in a list of values
func (f *MongoDBFiltering) NotIn(key string, values interface{}) interface{} {
	var valueList []string
	err := json.Unmarshal([]byte(values.(string)), &valueList)
	if err != nil {
		log.Errorf("Failed to unmarshal value list for NOT IN operation: %v", err)
		return nil
	}
	return bson.DocElem{
		Name: key,
		Value: bson.D{
			{Name: "$nin", Value: valueList},
		},
	}
}

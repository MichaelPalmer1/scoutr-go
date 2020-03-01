package azure

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/MichaelPalmer1/simple-api-go/providers/base"
	"github.com/globalsign/mgo/bson"
)

type MongoDBFiltering struct {
	base.Filtering
}

// Operations : Map of supported operations for this filter provider
func (f *MongoDBFiltering) Operations() base.OperationMap {
	return base.OperationMap{
		"ne":      f.NotEquals,
		"gt":      f.GreaterThan,
		"lt":      f.LessThan,
		"gte":     f.GreaterThanEqual,
		"lte":     f.LessThanEqual,
		"between": f.Between,
		//"in": f.In,
	}
}

// And : Takes two conditions and performs an AND operation on them
func (f *MongoDBFiltering) And(conditions, condition interface{}) interface{} {
	var output bson.D

	if _, ok := conditions.(bson.DocElem); ok {
		output = append(output, conditions.(bson.DocElem))
	} else {
		output = append(output, conditions.(bson.D)...)
	}

	if _, ok := condition.(bson.DocElem); ok {
		output = append(output, condition.(bson.DocElem))
	} else {
		output = append(output, condition.(bson.D)...)
	}

	return output
}

// Equals : Standard equals operation
func (f *MongoDBFiltering) Equals(key string, value interface{}) interface{} {
	return bson.DocElem{Name: key, Value: value}
}

// NotEquals : Standard not equals operation
func (f *MongoDBFiltering) NotEquals(key string, value interface{}) interface{} {
	return bson.DocElem{
		Name: key,
		Value: bson.D{{
			Name:  "$ne",
			Value: value,
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

// TODO: Find a way to use bson.A in github.com/globalsign/mgo/bson
// It is supported in the mongodb BSON, but not the globalsign one. If I use the mongodb
// one, the query breaks.

// In : Find all records with a list of values
//func (f *MongoDBFiltering) In(key string, values interface{}) interface{} {
//	return bson.E{
//		Key: key,
//		Value: bson.D{
//			{Key: "$in", Value: bson.A{values}},
//		},
//	}
//}

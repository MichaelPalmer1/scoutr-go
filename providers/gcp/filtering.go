package gcp

import (
	"encoding/json"

	"cloud.google.com/go/firestore"
	"github.com/MichaelPalmer1/scoutr-go/providers/base"
	log "github.com/sirupsen/logrus"
)

type FirestoreFiltering struct {
	base.BaseFiltering
	Query firestore.Query
}

// Operations : Map of supported operations for this filter provider
func (f *FirestoreFiltering) Operations() base.OperationMap {
	return base.OperationMap{
		"gt":      f.GreaterThan,
		"lt":      f.LessThan,
		"ge":      f.GreaterThanEqual,
		"le":      f.LessThanEqual,
		"between": f.Between,
		"in":      f.In,
	}
}

// And : Takes two conditions and performs an AND operation on them
func (f *FirestoreFiltering) And(conditions, condition interface{}) interface{} {
	// This is a bit hacky, but the best solution I have for now...
	return f.Query
}

// Equals : Standard equals operation
func (f *FirestoreFiltering) Equals(key string, value interface{}) interface{} {
	f.Query = f.Query.Where(key, "==", value)
	return f.Query
}

// GreaterThan : Check if a value is greater than a string
func (f *FirestoreFiltering) GreaterThan(key string, value interface{}) interface{} {
	f.Query = f.Query.Where(key, ">", value)
	return f.Query
}

// LessThan : Check if a value is greater than a string
func (f *FirestoreFiltering) LessThan(key string, value interface{}) interface{} {
	f.Query = f.Query.Where(key, "<", value)
	return f.Query
}

// GreaterThanEqual : Check if a value is greater than a string
func (f *FirestoreFiltering) GreaterThanEqual(key string, value interface{}) interface{} {
	f.Query = f.Query.Where(key, ">=", value)
	return f.Query
}

// LessThanEqual : Check if a value is greater than a string
func (f *FirestoreFiltering) LessThanEqual(key string, value interface{}) interface{} {
	f.Query = f.Query.Where(key, "<=", value)
	return f.Query
}

// Between : Check for records that are between a low and high value
func (f *FirestoreFiltering) Between(key string, value interface{}) interface{} {
	var valueList []string
	err := json.Unmarshal([]byte(value.(string)), &valueList)
	if err != nil {
		log.Errorf("Failed to unmarshal value list for BETWEEN operation: %v", err)
		return nil
	}
	f.Query = f.Query.Where(key, ">=", valueList[0]).Where(key, "<=", valueList[1])
	return f.Query
}

// In : Find all records with a list of values
func (f *FirestoreFiltering) In(key string, values interface{}) interface{} {
	var valueList []string
	err := json.Unmarshal([]byte(values.(string)), &valueList)
	if err != nil {
		log.Errorf("Failed to unmarshal value list for IN operation: %v", err)
		return nil
	}
	f.Query = f.Query.Where(key, "in", valueList)
	return f.Query
}

package gcp

import (
	"encoding/json"
	"reflect"
	"regexp"

	"cloud.google.com/go/firestore"
	"github.com/MichaelPalmer1/simple-api-go/models"
	log "github.com/sirupsen/logrus"
)

func buildFilters(user *models.User, filters map[string]string, collection *firestore.CollectionRef) (firestore.Query, error) {
	query := collection.Query
	re := regexp.MustCompile(`^(.+)__(in|contains|notcontains|startswith|ne|gt|lt|ge|le|between|exists)$`)

	// Build user filters
	if user != nil {
		for _, item := range user.FilterFields {
			if value, ok := item.Value.(string); ok {
				// Value is a single string
				query = query.Where(item.Field, "==", value)
			} else if value, ok := item.Value.([]interface{}); ok {
				// Value is a list of strings
				query = query.Where(item.Field, "array-contains-any", value)
			} else {
				log.Warnln("Received value of unknown type", item.Value)
				log.Warnln("Type", reflect.TypeOf(item.Value))
				continue
			}
		}
	}

	// Build specified filters
	for key, value := range filters {
		// Check for magic operator matches
		matches := re.FindAllStringSubmatch(key, -1)
		if len(matches) > 0 && len(matches[0]) == 3 {
			key = matches[0][1]
			operation := matches[0][2]

			// Perform filter based on the desired magic operator
			switch operation {
			case "in":
				var valueList []string
				json.Unmarshal([]byte(value), &valueList)
				query = query.Where(key, "in", valueList)
			// case "contains":
			// 	condition = attr.Contains(value)
			// case "notcontains":
			// 	condition = expression.Not(attr.Contains(value))
			// case "exists":
			// 	if value == "true" {
			// 		condition = attr.AttributeExists()
			// 	} else if value == "false" {
			// 		condition = attr.AttributeNotExists()
			// 	} else {
			// 		continue
			// 	}
			// case "startswith":
			// 	condition = attr.BeginsWith(value)
			// case "ne":
			// 	condition = attr.NotEqual(expression.Value(value))
			case "between":
				var valueList []string
				json.Unmarshal([]byte(value), &valueList)
				query = query.Where(key, ">=", valueList[0]).Where(key, "<=", valueList[1])
			case "gt":
				query = query.Where(key, ">", value)
			case "lt":
				query = query.Where(key, "<", value)
			case "ge":
				query = query.Where(key, ">=", value)
			case "le":
				query = query.Where(key, "<=", value)
			default:
				return query, &models.BadRequest{
					Message: "Unsupported magic operator",
				}
			}
		} else {
			query = query.Where(key, "==", value)
		}
	}

	return query, nil
}

func multiFilter(user *models.User, collection *firestore.CollectionRef, key string, values []string) (firestore.Query, error) {
	query, err := buildFilters(user, nil, collection)
	if err != nil {
		return collection.Query, err
	}

	// Build the condition
	// TODO: Handle if there are more than 10 values
	query = query.Where(key, "in", values)

	return query, nil
}

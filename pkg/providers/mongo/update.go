package mongo

// import (
// 	"fmt"

// 	"github.com/globalsign/mgo/bson"

// 	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
// 	log "github.com/sirupsen/logrus"
// )

// // Update : Update an item
// func (api MongoAPI) Update(request types.Request, partitionKey map[string]interface{}, item map[string]string, validation map[string]types.FieldValidation, auditAction string) (interface{}, error) {
// 	var output interface{}
// 	collection := api.Client.C(api.Config.DataTable)

// 	// Get the user
// 	user, err := api.InitializeRequest(request)
// 	if err != nil {
// 		// Bad user - pass the error through
// 		return nil, err
// 	}

// 	// Run data validation
// 	if validation != nil {
// 		err := api.ValidateFields(validation, item, nil, true)
// 		if err != nil {
// 			log.Errorln("Field validation error", err)
// 			return nil, err
// 		}
// 	}

// 	// Make sure user is not trying to update restricted fields
// 	for _, field := range user.UpdateFieldsRestricted {
// 		if _, ok := item[field]; ok {
// 			return nil, &types.BadRequest{
// 				Message: fmt.Sprintf("Not authorized to update field '%s'", field),
// 			}
// 		}
// 	}

// 	// Check all keys of the update input
// 	for _, key := range item {
// 		// Make sure fields being updated are not excluded from user
// 		for _, field := range user.ExcludeFields {
// 			if field == key {
// 				return nil, &types.BadRequest{
// 					Message: fmt.Sprintf("Not authorized to update field '%s'", key),
// 				}
// 			}
// 		}

// 		// Make sure the user is only updating permitted fields
// 		if len(user.UpdateFieldsPermitted) > 0 {
// 			found := false
// 			for _, field := range user.UpdateFieldsPermitted {
// 				if field == key {
// 					found = true
// 					break
// 				}
// 			}

// 			if !found {
// 				return nil, &types.BadRequest{
// 					Message: fmt.Sprintf("Not authorized to update field '%s'", key),
// 				}
// 			}
// 		}
// 	}

// 	// Build update expression
// 	var updates bson.D
// 	for key, value := range item {
// 		updates = append(updates, bson.DocElem{
// 			Name: "$set",
// 			Value: bson.D{{
// 				Name:  key,
// 				Value: value,
// 			}},
// 		})
// 	}

// 	// Build pre-condition filters. This will apply all the filter criteria for the user to this selector query and
// 	// throw a not found error if the user is not permitted to view the item
// 	rawFilters, _, err := api.Filter(
// 		&api.Filtering,
// 		user,
// 		map[string]string{api.Config.PrimaryKey: partitionKey[api.Config.PrimaryKey]},
// 	)
// 	if err != nil {
// 		log.Errorf("Error generating rawFilters: %v", err)
// 		return nil, err
// 	}

// 	// Make sure filters are cast as bson.D
// 	var selector bson.D
// 	if _, ok := rawFilters.(bson.DocElem); ok {
// 		// Single filter
// 		selector = bson.D{rawFilters.(bson.DocElem)}
// 	} else {
// 		// Multiple filters
// 		selector = rawFilters.(bson.D)
// 	}

// 	// Update the item
// 	err = collection.Update(selector, updates)
// 	if err != nil {
// 		if err.Error() == "not found" {
// 			return nil, &types.NotFound{
// 				Message: "Item does not exist or you do not have permission to view it",
// 			}
// 		} else {
// 			log.Errorln("Error while attempting to update item", err)
// 			return nil, err
// 		}
// 	}

// 	// Pull the updated item to show as a result
// 	err = collection.Find(bson.M{api.Config.PrimaryKey: partitionKey[api.Config.PrimaryKey]}).One(&output)
// 	if err != nil {
// 		log.Errorln("Failed to fetch updated item to show in results")
// 	}

// 	// Create audit log
// 	if err := api.auditLog(auditAction, request, *user, &partitionKey, &item); err != nil {
// 		log.Warnf("Failed to create audit log: %v", err)
// 	}

// 	return output, nil
// }

package gcp

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	log "github.com/sirupsen/logrus"
)

// Update : Update an item
func (api FirestoreAPI) Update(req models.Request, partitionKey map[string]string, item map[string]string, validation map[string]utils.FieldValidation, auditAction string) (interface{}, error) {
	var output interface{}

	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Run data validation
	if validation != nil {
		err := utils.ValidateFields(validation, item, nil, true)
		if err != nil {
			log.Errorln("Field validation error", err)
			return nil, err
		}
	}

	// Make sure user is not trying to update restricted fields
	for _, field := range user.UpdateFieldsRestricted {
		if _, ok := item[field]; ok {
			return nil, &models.BadRequest{
				Message: fmt.Sprintf("Not authorized to update field '%s'", field),
			}
		}
	}

	// Check all keys of the update input
	for key, _ := range item {
		// Make sure fields being updated are not excluded from user
		for _, field := range user.ExcludeFields {
			if field == key {
				return nil, &models.BadRequest{
					Message: fmt.Sprintf("Not authorized to update field '%s'", key),
				}
			}
		}

		// Make sure the user is only updating permitted fields
		if len(user.UpdateFieldsPermitted) > 0 {
			found := false
			for _, field := range user.UpdateFieldsPermitted {
				if field == key {
					found = true
					break
				}
			}

			if !found {
				return nil, &models.BadRequest{
					Message: fmt.Sprintf("Not authorized to update field '%s'", key),
				}
			}
		}
	}

	// Build update expression
	updates := []firestore.Update{}
	for key, value := range item {
		updates = append(updates, firestore.Update{
			Path:  key,
			Value: value,
		})
	}

	// Update the item
	collection := api.Client.Collection(api.Config.DataTable)
	_, err = collection.Doc(partitionKey[api.Config.PrimaryKey]).Update(api.context, updates)
	if err != nil {
		log.Errorln("Error while attempting to update item", err)
		return nil, err
	}

	// Pull the updated item to show as a result
	doc, err := collection.Doc(partitionKey[api.Config.PrimaryKey]).Get(api.context)
	if err != nil {
		log.Errorln("Failed to fetch updated item to show in results")
	} else {
		output = doc.Data()
	}

	// Create audit log
	api.auditLog("UPDATE", req, *user, &partitionKey, &item)

	return output, nil
}

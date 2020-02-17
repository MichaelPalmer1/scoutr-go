package gcp

import (
	"cloud.google.com/go/firestore"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	log "github.com/sirupsen/logrus"
)

// Update : Update an item
func (api *FirestoreAPI) Update(req models.Request, partitionKey map[string]string, item map[string]string, validation map[string]utils.FieldValidation, auditAction string) (interface{}, error) {
	var output interface{}

	// Get the user
	_, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Run data validation
	if validation != nil {
		log.Infoln("Running field validation")
		err := utils.ValidateFields(validation, item, nil, true)
		if err != nil {
			log.Errorln("Field validation error", err)
			return nil, err
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

	// TODO: Build pre-condition filters

	// Update the item
	collection := api.client.Collection(api.Config.DataTable)
	_, err = collection.Doc(partitionKey[api.Config.PrimaryKey]).Update(api.context, updates)
	if err != nil {
		log.Errorln("Error while attempting to update item", err)

		// Check if this was a conditional check failure
		// if _, ok := err.(*dynamodb.ConditionalCheckFailedException); ok {
		// 	return nil, &models.BadRequest{
		// 		Message: "Item does not exist or you do not have permission to update it",
		// 	}
		// }

		return nil, err
	}

	// Unmarshal into output interface
	// dynamodbattribute.UnmarshalMap(updatedItem.Attributes, &output)

	// Create audit log
	// api.auditLog("UPDATE", req, *user, &partitionKey, &item)

	return output, nil
}

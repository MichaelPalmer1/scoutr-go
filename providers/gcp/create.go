package gcp

import (
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	log "github.com/sirupsen/logrus"
)

// Create : Create an item
func (api *FirestoreAPI) Create(req models.Request, item map[string]string, validation map[string]utils.FieldValidation) error {
	// Get the user
	_, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return err
	}

	// Run data validation
	if validation != nil {
		log.Infoln("Running field validation")
		err := utils.ValidateFields(validation, item, nil, false)
		if err != nil {
			log.Errorln("Field validation error", err)
			return err
		}
	}

	collection := api.client.Collection(api.Config.DataTable)
	doc := collection.Doc(item[api.Config.PrimaryKey])

	// TODO: Build pre-condition filters. This may take more work with Firestore...

	_, err = doc.Create(api.context, item)
	if err != nil {
		log.Errorln("Error while attempting to create item", err)

		// Check if this was a conditional check failure
		// if _, ok := err.(*dynamodb.ConditionalCheckFailedException); ok {
		// 	return &models.BadRequest{
		// 		Message: "Item already exists or you do not have permission to create it",
		// 	}
		// }

		return err
	}

	// Create audit log
	// api.auditLog("CREATE", req, *user, &map[string]string{partitionKey: item[partitionKey]}, nil)

	return nil
}

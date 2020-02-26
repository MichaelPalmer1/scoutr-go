package gcp

import (
	"github.com/MichaelPalmer1/simple-api-go/models"
	log "github.com/sirupsen/logrus"
)

// Delete : Delete an item
func (api FirestoreAPI) Delete(req models.Request, partitionKey map[string]string) error {
	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return err
	}

	// Attempt to fetch the item
	if _, err := api.fetchItem(user, partitionKey[api.Config.PrimaryKey]); err != nil {
		return err
	}

	// Delete the item
	doc := api.Client.Collection(api.Config.DataTable).Doc(partitionKey[api.Config.PrimaryKey])
	_, err = doc.Delete(api.context)
	if err != nil {
		log.Errorln("Error while attempting to delete item", err)
		return err
	}

	// Create audit log
	api.auditLog("DELETE", req, *user, &partitionKey, nil)

	return nil
}

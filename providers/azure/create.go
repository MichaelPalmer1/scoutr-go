package azure

import (
	"github.com/MichaelPalmer1/scoutr-go/models"
	"github.com/MichaelPalmer1/scoutr-go/utils"
	log "github.com/sirupsen/logrus"
)

// Create : Create an item
func (api MongoDBAPI) Create(req models.Request, item map[string]string, validation map[string]utils.FieldValidation) error {
	// Get the user
	user, err := api.InitializeRequest(api, req)
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

	// TODO: Build pre-condition filters

	// Insert the item
	collection := api.Client.C(api.Config.DataTable)
	err = collection.Insert(item)
	if err != nil {
		return err
	}

	// Create audit log
	partitionKey := map[string]string{api.Config.PrimaryKey: item[api.Config.PrimaryKey]}
	if err := api.auditLog("CREATE", req, *user, &partitionKey, nil); err != nil {
		log.Warnf("Failed to create audit log: %v", err)
	}

	return nil
}

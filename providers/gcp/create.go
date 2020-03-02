package gcp

import (
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create : Create an item
func (api FirestoreAPI) Create(req models.Request, item map[string]string, validation map[string]utils.FieldValidation) error {
	// Get the user
	user, err := api.InitializeRequest(req)
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

	collection := api.Client.Collection(api.Config.DataTable)
	doc := collection.Doc(item[api.Config.PrimaryKey])

	// TODO: Build pre-condition filters. This may take more work with Firestore...

	_, err = doc.Create(api.context, item)
	if err != nil {
		// Attempt to convert error to a status code
		code, ok := status.FromError(err)

		// Check if the status conversion was successful
		if ok {
			switch code.Code() {
			case codes.AlreadyExists:
				return &models.BadRequest{
					Message: "Item already exists",
				}
			case codes.InvalidArgument:
				// Return bad request on invalid argument errors
				return &models.BadRequest{
					Message: code.Message(),
				}
			}
		}

		// Fallback to just returning the raw error
		return err
	}

	// Create audit log
	partitionKey := map[string]string{api.Config.PrimaryKey: doc.ID}
	if err := api.auditLog("CREATE", req, *user, &partitionKey, nil); err != nil {
		log.Warnf("Failed to create audit log: %v", err)
	}

	return nil
}

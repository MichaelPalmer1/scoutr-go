package mongo

// import (
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
// 	log "github.com/sirupsen/logrus"
// )

// // Create : Create an item
// func (api MongoAPI) Create(request types.Request, item map[string]interface{}, validation map[string]types.FieldValidation, requiredFields []string) error {
// 	// Get the user
// 	user, err := api.InitializeRequest(request)
// 	if err != nil {
// 		// Bad user - pass the error through
// 		return err
// 	}

// 	// Run data validation
// 	if validation != nil {
// 		log.Infoln("Running field validation")
// 		err := api.ValidateFields(validation, requiredFields, item, nil)
// 		if err != nil {
// 			log.Errorln("Field validation error", err)
// 			return err
// 		}
// 	}

// 	// TODO: Build pre-condition filters

// 	// Insert the item
// 	collection := api.Client.C(api.Config.DataTable)
// 	err = collection.Insert(item)
// 	if err != nil {
// 		return err
// 	}

// 	// Create audit log
// 	partitionKey := map[string]interface{}{api.Config.PrimaryKey: item[api.Config.PrimaryKey]}
// 	if err := api.auditLog(base.AuditActionCreate, request, *user, &partitionKey, nil); err != nil {
// 		log.Warnf("Failed to create audit log: %v", err)
// 	}

// 	return nil
// }

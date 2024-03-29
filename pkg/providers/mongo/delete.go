package mongo

// import (
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
// 	"github.com/globalsign/mgo/bson"
// 	log "github.com/sirupsen/logrus"
// )

// // Delete : Delete an item
// func (api MongoAPI) Delete(req types.Request, partitionKey map[string]interface{}) error {
// 	// Get the user
// 	user, err := api.InitializeRequest(req)
// 	if err != nil {
// 		// Bad user - pass the error through
// 		return err
// 	}

// 	// Build pre-condition filters. This will apply all the filter criteria for the user to this selector query and
// 	// throw a not found error if the user is not permitted to view the item
// 	rawFilters, err := api.Filtering.Filter(
// 		user,
// 		map[string][]string{api.Config.PrimaryKey: []string{partitionKey[api.Config.PrimaryKey].(string)}},
// 		base.FilterActionDelete,
// 	)
// 	if err != nil {
// 		log.Errorf("Error generating rawFilters: %v", err)
// 		return err
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

// 	// Delete the item
// 	collection := api.Client.C(api.Config.DataTable)
// 	err = collection.Remove(selector)
// 	if err != nil {
// 		if err.Error() == "not found" {
// 			return &types.NotFound{
// 				Message: "Item does not exist or you do not have permission to view it",
// 			}
// 		} else {
// 			log.Errorln("Error while attempting to delete item", err)
// 			return err
// 		}
// 	}

// 	// Create audit log
// 	if err := api.auditLog(base.AuditActionDelete, req, *user, &partitionKey, nil); err != nil {
// 		log.Warnf("Failed to create audit log: %v", err)
// 	}

// 	return nil
// }

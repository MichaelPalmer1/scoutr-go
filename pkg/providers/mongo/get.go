package mongo

// import (
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
// 	"github.com/globalsign/mgo/bson"
// 	log "github.com/sirupsen/logrus"
// )

// // Get : Get an item from the table
// func (api MongoAPI) Get(req types.Request, id string) (types.Record, error) {
// 	var record types.Record
// 	collection := api.Client.C(api.Config.DataTable)

// 	// Get the user
// 	user, err := api.InitializeRequest(req)
// 	if err != nil {
// 		// Bad user - pass the error through
// 		return nil, err
// 	}

// 	// Build pre-condition filters. This will apply all the filter criteria for the user to this selector query and
// 	// throw a not found error if the user is not permitted to view the item
// 	rawFilters, err := api.Filtering.Filter(
// 		user,
// 		map[string][]string{api.Config.PrimaryKey: []string{id}},
// 		base.FilterActionRead,
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

// 	// Fetch the item
// 	err = collection.Find(selector).One(&record)
// 	if err != nil {
// 		if err.Error() == "not found" {
// 			return nil, &types.NotFound{
// 				Message: "Item does not exist or you do not have permission to view it",
// 			}
// 		} else {
// 			// Pass through any errors
// 			return nil, err
// 		}
// 	}

// 	// Create audit log
// 	partitionKey := map[string]interface{}{api.Config.PrimaryKey: id}
// 	if err := api.auditLog(base.AuditActionGet, req, *user, &partitionKey, nil); err != nil {
// 		log.Warnf("Failed to create audit log: %v", err)
// 	}

// 	return record, nil
// }

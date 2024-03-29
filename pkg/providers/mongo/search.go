package mongo

// import (
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
// 	"github.com/globalsign/mgo/bson"
// 	log "github.com/sirupsen/logrus"
// )

// // Search : Search items in the table
// func (api MongoAPI) Search(req types.Request, key string, values []string) ([]types.Record, error) {
// 	records := []types.Record{}
// 	collection := api.Client.C(api.Config.DataTable)

// 	// Get the user
// 	user, err := api.InitializeRequest(req)
// 	if err != nil {
// 		// Bad user - pass the error through
// 		return nil, err
// 	}

// 	// Build filters
// 	rawFilters, err := api.Filtering.MultiFilter(user, key, values)
// 	if err != nil {
// 		log.Errorln("Error encountered during filtering", err)
// 		return nil, err
// 	}

// 	// Make sure filters are cast as bson.D
// 	var filters bson.D
// 	if rawFilters == nil {
// 		// No filters
// 		filters = nil
// 	} else if _, ok := rawFilters.(bson.DocElem); ok {
// 		// Single filter
// 		filters = bson.D{rawFilters.(bson.DocElem)}
// 	} else {
// 		// Multiple filters
// 		filters = rawFilters.(bson.D)
// 	}

// 	// Query the data
// 	err = collection.Find(filters).All(&records)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Filter the response
// 	api.PostProcess(records, user)

// 	// Create audit log
// 	if err := api.auditLog(base.AuditActionSearch, req, *user, nil, nil); err != nil {
// 		log.Warnf("Failed to create audit log: %v", err)
// 	}

// 	return records, nil
// }

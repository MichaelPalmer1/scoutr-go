package mongo

// import (
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
// 	"github.com/globalsign/mgo/bson"
// 	log "github.com/sirupsen/logrus"
// )

// // List : Lists all items in a table
// func (api MongoAPI) List(req types.Request) ([]types.Record, error) {
// 	records := []types.Record{}
// 	collection := api.Client.C(api.Config.DataTable)

// 	// Get the user
// 	user, err := api.InitializeRequest(req)
// 	if err != nil {
// 		// Bad user - pass the error through
// 		return nil, err
// 	}

// 	// Build rawFilters
// 	rawFilters, err := api.Filtering.Filter(user, api.BuildParams(req), base.FilterActionRead)
// 	if err != nil {
// 		log.Errorf("Error generating rawFilters: %v", err)
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
// 		log.Errorf("Error querying collection: %v", err)
// 		return nil, err
// 	}

// 	// Filter the response
// 	api.PostProcess(records, user)

// 	// Create audit log
// 	if err := api.auditLog(base.AuditActionList, req, *user, nil, nil); err != nil {
// 		log.Warnf("Failed to create audit log: %v", err)
// 	}

// 	return records, nil
// }

// // ListUniqueValues : Lists unique values for a column
// func (api MongoAPI) ListUniqueValues(req types.Request, uniqueKey string) ([]string, error) {
// 	records := []string{}
// 	collection := api.Client.C(api.Config.DataTable)

// 	// Get the user
// 	user, err := api.InitializeRequest(req)
// 	if err != nil {
// 		// Bad user - pass the error through
// 		return nil, err
// 	}

// 	// Build rawFilters
// 	rawFilters, err := api.Filtering.Filter(user, api.BuildParams(req), base.FilterActionRead)
// 	if err != nil {
// 		log.Errorf("Error generating rawFilters: %v", err)
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
// 	err = collection.Find(filters).Distinct(uniqueKey, &records)
// 	if err != nil {
// 		log.Errorf("Error querying collection: %v", err)
// 		return nil, err
// 	}

// 	// TODO: Filter the response
// 	//api.PostProcess(records, user)

// 	// Create audit log
// 	if err := api.auditLog(base.AuditActionList, req, *user, nil, nil); err != nil {
// 		log.Warnf("Failed to create audit log: %v", err)
// 	}

// 	return records, nil
// }

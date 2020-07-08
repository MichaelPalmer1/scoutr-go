package azure

import (
	"github.com/MichaelPalmer1/scoutr-go/models"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

// List : Lists all items in a table
func (api MongoDBAPI) List(req models.Request) ([]models.Record, error) {
	records := []models.Record{}
	collection := api.Client.C(api.Config.DataTable)

	// Get the user
	user, err := api.InitializeRequest(api, req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Build rawFilters
	rawFilters, _, err := api.Filter(&api.Filtering, user, api.BuildParams(req))
	if err != nil {
		log.Errorf("Error generating rawFilters: %v", err)
		return nil, err
	}

	// Make sure filters are cast as bson.D
	var filters bson.D
	if rawFilters == nil {
		// No filters
		filters = nil
	} else if _, ok := rawFilters.(bson.DocElem); ok {
		// Single filter
		filters = bson.D{rawFilters.(bson.DocElem)}
	} else {
		// Multiple filters
		filters = rawFilters.(bson.D)
	}

	// Query the data
	err = collection.Find(filters).All(&records)
	if err != nil {
		log.Errorf("Error querying collection: %v", err)
		return nil, err
	}

	// Filter the response
	api.PostProcess(records, user)

	// Create audit log
	if err := api.auditLog("LIST", req, *user, nil, nil); err != nil {
		log.Warnf("Failed to create audit log: %v", err)
	}

	return records, nil
}

// ListUniqueValues : Lists unique values for a column
func (api MongoDBAPI) ListUniqueValues(req models.Request, uniqueKey string) ([]string, error) {
	records := []string{}
	collection := api.Client.C(api.Config.DataTable)

	// Get the user
	user, err := api.InitializeRequest(api, req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Build rawFilters
	rawFilters, _, err := api.Filter(&api.Filtering, user, api.BuildParams(req))
	if err != nil {
		log.Errorf("Error generating rawFilters: %v", err)
		return nil, err
	}

	// Make sure filters are cast as bson.D
	var filters bson.D
	if rawFilters == nil {
		// No filters
		filters = nil
	} else if _, ok := rawFilters.(bson.DocElem); ok {
		// Single filter
		filters = bson.D{rawFilters.(bson.DocElem)}
	} else {
		// Multiple filters
		filters = rawFilters.(bson.D)
	}

	// Query the data
	err = collection.Find(filters).Distinct(uniqueKey, &records)
	if err != nil {
		log.Errorf("Error querying collection: %v", err)
		return nil, err
	}

	// TODO: Filter the response
	//api.PostProcess(records, user)

	// Create audit log
	if err := api.auditLog("LIST", req, *user, nil, nil); err != nil {
		log.Warnf("Failed to create audit log: %v", err)
	}

	return records, nil
}

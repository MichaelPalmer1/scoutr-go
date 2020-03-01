package azure

import (
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

func (api MongoDBAPI) List(req models.Request) ([]models.Record, error) {
	records := []models.Record{}
	collection := api.Client.C(api.Config.DataTable)

	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Copy queryParams into params
	params := make(map[string]string)
	for key, value := range req.QueryParams {
		params[key] = value
	}

	// Merge pathParams into params
	for key, value := range req.PathParams {
		params[key] = value
	}

	// Generate dynamic search
	searchKey, hasSearchKey := req.PathParams["search_key"]
	searchValue, hasSearchValue := req.PathParams["search_value"]
	if hasSearchKey && hasSearchValue {
		// Map the search key and value into path params
		params[searchKey] = searchValue
		delete(params, "search_key")
		delete(params, "search_value")
	}

	// Build rawFilters
	rawFilters, _, err := api.Filter(&api.Filtering, user, params)
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

package gcp

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/providers/base"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

// FirestoreAPI : API, based off of SimpleAPI, used to talk to Google Firestore
type FirestoreAPI struct {
	*base.SimpleAPI
	client  *firestore.Client
	context context.Context
}

// Init : Initialize the Firestore client
func (api *FirestoreAPI) Init(projectID string, options option.ClientOption) {
	api.context = context.Background()
	client, err := firestore.NewClient(api.context, projectID, options)
	if err != nil {
		log.Fatalln("Failed to initialize Firestore client", err)
	}
	api.client = client
}

// List : List all records
func (api *FirestoreAPI) List(req models.Request) ([]models.Record, error) {
	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	log.Infof("Found user: %v", user)

	// Generate dynamic search
	searchKey, hasSearchKey := req.PathParams["search_key"]
	searchValue, hasSearchValue := req.PathParams["search_value"]
	if hasSearchKey && hasSearchValue {
		// Map the search key and value into path params
		req.PathParams[searchKey] = searchValue
		delete(req.PathParams, "search_key")
		delete(req.PathParams, "search_value")
	}

	// Merge pathParams into queryParams
	for key, value := range req.PathParams {
		req.QueryParams[key] = value
	}

	// Build filters
	collection, err := buildFilters(user, req.QueryParams, api.client.Collection(api.Config.DataTable))
	if err != nil {
		return nil, err
	}

	// Download the data
	docs, err := collection.Documents(api.context).GetAll()
	if err != nil {
		log.Errorln("Error while attempting to list records", err)
		return nil, nil
	}

	records := []models.Record{}
	for _, doc := range docs {
		records = append(records, doc.Data())
	}

	// Filter the response
	api.PostProcess(records, user)

	return records, nil
}

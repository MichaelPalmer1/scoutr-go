package gcp

import (
	"cloud.google.com/go/firestore"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// List : List all records
func (api FirestoreAPI) List(req models.Request) ([]models.Record, error) {
	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

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
	collection := api.Client.Collection(api.Config.DataTable)
	f := FirestoreFiltering{
		Query: collection.Query,
	}
	filters, _, err := api.Filter(&f, user, req.QueryParams)
	if err != nil {
		return nil, err
	}
	query := collection.Query
	if filters != nil {
		query = filters.(firestore.Query)
	}

	// Download the data
	docs, err := query.Documents(api.context).GetAll()
	if err != nil {
		// Attempt to convert error to a status code
		code, ok := status.FromError(err)

		// Check if the status conversion was successful
		if ok {
			switch code.Code() {
			case codes.InvalidArgument:
				// Return bad request on invalid argument errors
				return nil, &models.BadRequest{
					Message: code.Message(),
				}
			}
		}

		// Fallback to just returning the raw error
		return nil, err
	}

	// TODO: fix this, this feels hacky...and not optimal
	records := []models.Record{}
	for _, doc := range docs {
		records = append(records, doc.Data())
	}

	// Filter the response
	api.PostProcess(records, user)

	return records, nil
}

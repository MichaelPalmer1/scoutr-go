package gcp

import (
	"cloud.google.com/go/firestore"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get : Get an item from the table
func (api FirestoreAPI) Get(req models.Request, id string) (models.Record, error) {
	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Build filters
	collection := api.Client.Collection(api.Config.DataTable)
	f := FirestoreFiltering{
		Query: collection.Query,
	}
	filters, _, err := api.Filter(&f, user, map[string]string{})
	if err != nil {
		return nil, err
	}
	query := collection.Query
	if filters != nil {
		query = filters.(firestore.Query)
	}

	// Build key condition
	query = query.Where(api.Config.PrimaryKey, "==", id)

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

	// Make sure only a single record was returned
	if len(records) > 1 {
		return nil, &models.BadRequest{
			Message: "Multiple items returned",
		}
	} else if len(records) == 0 {
		return nil, &models.NotFound{
			Message: "Item does not exist or you do not have permission to view it",
		}
	}

	// Create audit log
	api.auditLog("GET", req, *user, &map[string]string{api.Config.PrimaryKey: id}, nil)

	return records[0], nil
}

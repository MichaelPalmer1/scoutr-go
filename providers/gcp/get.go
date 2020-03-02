package gcp

import (
	"cloud.google.com/go/firestore"
	"github.com/MichaelPalmer1/simple-api-go/models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get : Get an item from the table
func (api FirestoreAPI) Get(req models.Request, id string) (models.Record, error) {
	// Get the user
	user, err := api.InitializeRequest(api, req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Fetch the item
	record, err := api.fetchItem(user, id)
	if err != nil {
		// Pass through any errors
		return nil, err
	}

	// Create audit log
	partitionKey := map[string]string{api.Config.PrimaryKey: id}
	if err := api.auditLog("GET", req, *user, &partitionKey, nil); err != nil {
		log.Warnf("Failed to create audit log: %v", err)
	}

	return record, nil
}

// Attempt to fetch a single item, applying any user filters beforehand
func (api *FirestoreAPI) fetchItem(user *models.User, id string) (models.Record, error) {
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

	// Query the data
	iter := query.Documents(api.context)
	records := []models.Record{}

	// Iterate through the results
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
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

		// Add item to records
		records = append(records, doc.Data())
	}

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

	return records[0], nil
}

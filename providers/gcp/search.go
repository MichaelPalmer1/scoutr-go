package gcp

import (
	"cloud.google.com/go/firestore"
	"github.com/MichaelPalmer1/simple-api-go/models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Search : Search items in the table
func (api FirestoreAPI) Search(req models.Request, key string, values []string) ([]models.Record, error) {
	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Build filters
	f := FirestoreFiltering{
		Query: api.Client.Collection(api.Config.DataTable).Query,
	}
	filters, err := api.MultiFilter(&f, user, key, values)
	if err != nil {
		log.Errorln("Error encountered during filtering", err)
		return nil, err
	}

	// Cast to Firestore query
	query := filters.(firestore.Query)

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

	// Filter the response
	api.PostProcess(records, user)

	// Create audit log
	api.auditLog("SEARCH", req, *user, nil, nil)

	return records, nil
}

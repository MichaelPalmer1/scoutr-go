package gcp

import (
	"cloud.google.com/go/firestore"
	"github.com/MichaelPalmer1/simple-api-go/models"
	log "github.com/sirupsen/logrus"
)

// Search : Search items in the table
func (api *FirestoreAPI) Search(req models.Request, key string, values []string) ([]models.Record, error) {
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

	// Download the data
	query := filters.(firestore.Query)
	data, err := query.Documents(api.context).GetAll()
	if err != nil {
		log.Errorln("Error while attempting to list records", err)
		return nil, nil
	}

	// TODO: fix this, this feels hacky...and not optimal
	records := []models.Record{}
	for _, item := range data {
		records = append(records, item.Data())
	}

	// Filter the response
	api.PostProcess(records, user)

	// Create audit log
	// api.auditLog("SEARCH", req, *user, nil, nil)

	return records, nil
}

package gcp

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/MichaelPalmer1/simple-api-go/providers/base"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

// FirestoreAPI : API, based off of SimpleAPI, used to talk to Google Firestore
type FirestoreAPI struct {
	*base.SimpleAPI
	Client  *firestore.Client
	context context.Context
}

// Init : Initialize the Firestore client
func (api *FirestoreAPI) Init(projectID string, options option.ClientOption) {
	api.context = context.Background()
	client, err := firestore.NewClient(api.context, projectID, options)
	if err != nil {
		log.Fatalln("Failed to initialize Firestore client", err)
	}
	api.Client = client
}

// Close : Close all connections
func (api *FirestoreAPI) Close() {
	api.Client.Close()
}

package main

import (
	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/providers/base"
	"github.com/MichaelPalmer1/simple-api-go/providers/gcp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func main() {
	request := models.Request{
		QueryParams: map[string]string{
			"value": `hello`,
			// "type__in": `[1,2,3,4,5,6,7,8,9,10]`,
		},
		PathParams: map[string]string{},
		Method:     "GET",
		Path:       "/items",
		UserAgent:  "test",
		SourceIP:   "1.2.3.4",
		User: models.RequestUser{
			ID: "fff",
			Data: &models.UserData{
				Username: "mp",
				Name:     "Michael",
				Email:    "michael@michael.michael",
				Groups:   []string{"michael2"},
			},
		},
	}

	gcp := gcp.FirestoreAPI{
		SimpleAPI: &base.SimpleAPI{
			Config: config.Config{
				DataTable:  "data",
				AuthTable:  "auth",
				AuditTable: "audit",
				GroupTable: "groups",
			},
		},
	}

	// Initialize the client
	gcp.Init("simple-api-265401", option.WithCredentialsFile("/home/michael/Downloads/gcp.json"))
	defer gcp.Close()

	// List the records
	records, err := gcp.List(request)
	if err != nil {
		log.Errorln(err)
		return
	}

	// Show results
	for _, doc := range records {
		log.Infoln(doc)
	}
}

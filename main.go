package main

import (
	"os/user"
	"path/filepath"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/models"
	dynamo "github.com/MichaelPalmer1/simple-api-go/providers/aws"
	"github.com/MichaelPalmer1/simple-api-go/providers/base"
	"github.com/MichaelPalmer1/simple-api-go/providers/gcp"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

var request models.Request

func google() {
	api := gcp.FirestoreAPI{
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
	api.Init("simple-api-265401", option.WithCredentialsFile("/home/michael/Downloads/gcp.json"))
	defer api.Close()

	// List the records
	records, err := api.List(request)
	if err != nil {
		log.Errorln(err)
		return
	}

	// Show results
	for _, doc := range records {
		log.Infoln(doc)
	}
}

func amazon() {
	api := dynamo.DynamoAPI{
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
	usr, _ := user.Current()
	creds := credentials.NewSharedCredentials(filepath.Join(usr.HomeDir, ".aws/credentials"), "default")
	api.Init(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: creds,
	})

	// List the records
	records, err := api.List(request)
	if err != nil {
		log.Errorln(err)
		return
	}

	// Show results
	for _, doc := range records {
		log.Infoln(doc)
	}
}

func init() {
	request = models.Request{
		QueryParams: map[string]string{
			"value": `hello`,
			// "type__ge":  "test",
			// "name__in":  `["ff6", "f"]`,
			"value__gt": "a",
			"value__lt": "z",
			// "key__ge":        "t",
			"value__le":      "v",
			"value__between": `["c", "m"]`,
		},
		PathParams: map[string]string{},
		Method:     "GET",
		Path:       "/items",
		UserAgent:  "test",
		SourceIP:   "1.2.3.4",
		User: models.RequestUser{
			ID: "michael",
			// Data: &models.UserData{
			// 	Username: "mp",
			// 	Name:     "Michael",
			// 	Email:    "michael@michael.michael",
			// 	Groups:   []string{"michael2"},
			// },
		},
	}
}

func main() {
	google()
	// amazon()
}

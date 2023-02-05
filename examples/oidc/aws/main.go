package main

import (
	"flag"
	"net/http"
	"strings"

	"github.com/MichaelPalmer1/scoutr-go/pkg/config"
	"github.com/MichaelPalmer1/scoutr-go/pkg/helpers"
	dynamo "github.com/MichaelPalmer1/scoutr-go/pkg/providers/aws"
	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	"github.com/aws/aws-sdk-go-v2/aws"
	log "github.com/sirupsen/logrus"
)

func init() {
	api = dynamo.DynamoAPI{
		Scoutr: &base.Scoutr{
			Config: config.Config{},
		},
		//Filtering: dynamo.DynamoFiltering{},
	}
}

func main() {
	// Command line arguments
	var nameHeader string
	flag.StringVar(&api.Config.DataTable, "data-table", "", "Data table")
	flag.StringVar(&api.Config.AuthTable, "auth-table", "", "Auth table")
	flag.StringVar(&api.Config.GroupTable, "group-table", "", "Group table")
	flag.StringVar(&api.Config.AuditTable, "audit-table", "", "Audit table")
	flag.IntVar(&api.Config.LogRetentionDays, "log-retention-days", 30, "Days to retain read logs")
	flag.StringVar(&api.Config.OIDCUsernameHeader, "oidc-username-header", "Oidc-Claim-Sub", "Username header from OIDC")
	flag.StringVar(&nameHeader, "oidc-name-header", "Oidc-Claim-Name", "Name header from OIDC")
	flag.StringVar(&api.Config.OIDCEmailHeader, "oidc-email-header", "Oidc-Claim-Mail", "Email header from OIDC")
	flag.StringVar(&api.Config.OIDCGroupHeader, "oidc-group-header", "", "Group header from OIDC")
	flag.Parse()

	api.Config.OIDCNameHeader = strings.Split(nameHeader, ",")

	// Make sure required fields are provided
	if api.Config.DataTable == "" {
		log.Fatalln("data-table argument is required")
	}
	if api.Config.AuthTable == "" {
		log.Fatalln("auth-table argument is required")
	}
	if api.Config.GroupTable == "" {
		log.Fatalln("group-table argument is required")
	}

	config := aws.NewConfig()
	config.Region = "us-east-1"
	api.Init(*config)

	// Initialize http server
	router, err := helpers.InitHTTPServer(api, "/items/")
	if err != nil {
		panic(err)
	}

	// Add get/create/update endpoints
	router.POST("/item/", create)
	router.GET("/item/:id", get)
	router.PUT("/item/:id", update)
	router.DELETE("/item/:id", delete)
	router.GET("/types/", listTypes)

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", router))
}

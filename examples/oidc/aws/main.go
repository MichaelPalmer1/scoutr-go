package main

import (
	"flag"
	"net/http"
	"strings"

	scoutrConfig "github.com/MichaelPalmer1/scoutr-go/pkg/config"
	"github.com/MichaelPalmer1/scoutr-go/pkg/helpers"
	dynamo "github.com/MichaelPalmer1/scoutr-go/pkg/providers/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Command line arguments
	var nameHeader string
	var conf scoutrConfig.Config

	flag.StringVar(&conf.DataTable, "data-table", "", "Data table")
	flag.StringVar(&conf.AuthTable, "auth-table", "", "Auth table")
	flag.StringVar(&conf.GroupTable, "group-table", "", "Group table")
	flag.StringVar(&conf.AuditTable, "audit-table", "", "Audit table")
	flag.IntVar(&conf.LogRetentionDays, "log-retention-days", 30, "Days to retain read logs")
	flag.StringVar(&conf.OIDCUsernameHeader, "oidc-username-header", "Oidc-Claim-Sub", "Username header from OIDC")
	flag.StringVar(&nameHeader, "oidc-name-header", "Oidc-Claim-Name", "Name header from OIDC")
	flag.StringVar(&conf.OIDCEmailHeader, "oidc-email-header", "Oidc-Claim-Mail", "Email header from OIDC")
	flag.StringVar(&conf.OIDCGroupHeader, "oidc-group-header", "", "Group header from OIDC")
	flag.Parse()

	conf.OIDCNameHeader = strings.Split(nameHeader, ",")

	// Make sure required fields are provided
	if conf.DataTable == "" {
		log.Fatalln("data-table argument is required")
	}
	if conf.AuthTable == "" {
		log.Fatalln("auth-table argument is required")
	}
	if conf.GroupTable == "" {
		log.Fatalln("group-table argument is required")
	}

	awsConfig := aws.NewConfig()
	awsConfig.Region = "us-east-1"
	api := dynamo.NewDynamoAPI(conf, *awsConfig)

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

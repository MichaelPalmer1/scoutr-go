package main

import (
	"flag"
	"net/http"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/helpers"
	"github.com/MichaelPalmer1/simple-api-go/providers/azure"
	"github.com/MichaelPalmer1/simple-api-go/providers/base"
	log "github.com/sirupsen/logrus"
)

func init() {
	api = azure.MongoDBAPI{
		SimpleAPI: &base.SimpleAPI{
			Config: config.Config{
				PrimaryKey: "id",
			},
		},
	}
}

func main() {
	var address, database, username, password string
	// Command line arguments
	flag.StringVar(&api.Config.DataTable, "data-table", "", "Data table")
	flag.StringVar(&api.Config.AuthTable, "auth-table", "", "Auth table")
	flag.StringVar(&api.Config.GroupTable, "group-table", "", "Group table")
	flag.StringVar(&api.Config.AuditTable, "audit-table", "", "Audit table")
	flag.IntVar(&api.Config.LogRetentionDays, "log-retention-days", 30, "Days to retain read logs")
	flag.StringVar(&api.Config.OIDCUsernameHeader, "oidc-username-header", "Oidc-Claim-Sub", "Username header from OIDC")
	flag.StringVar(&api.Config.OIDCNameHeader, "oidc-name-header", "Oidc-Claim-Name", "Name header from OIDC")
	flag.StringVar(&api.Config.OIDCEmailHeader, "oidc-email-header", "Oidc-Claim-Mail", "Email header from OIDC")
	flag.StringVar(&api.Config.OIDCGroupHeader, "oidc-group-header", "", "Group header from OIDC")
	flag.StringVar(&address, "address", "", "MongoDB server address formatted as HOST:PORT")
	flag.StringVar(&database, "database", "", "Name of the database")
	flag.StringVar(&username, "username", "", "MongoDB username to authenticate with")
	flag.StringVar(&password, "password", "", "MongoDB password to authenticate with")
	flag.Parse()

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

	// Initialize the client
	api.Init(
		address,
		database,
		username,
		password,
	)
	defer api.Close()

	// Initialize http server
	router, err := helpers.InitHTTPServer(api, "id", "/items/", []string{"CREATE", "UPDATE"})
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

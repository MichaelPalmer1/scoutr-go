package main

import (
	"flag"
	"net/http"

	"github.com/MichaelPalmer1/scoutr-go/config"
	"github.com/MichaelPalmer1/scoutr-go/helpers"
	"github.com/MichaelPalmer1/scoutr-go/providers/base"
	"github.com/MichaelPalmer1/scoutr-go/providers/gcp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func init() {
	api = gcp.FirestoreAPI{
		Scoutr: &base.Scoutr{
			Config: config.Config{
				PrimaryKey: "id",
			},
		},
	}
}

func main() {
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
	api.Init("scoutr-1234", option.WithCredentialsFile("/home/scoutr/gcp.json"))
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

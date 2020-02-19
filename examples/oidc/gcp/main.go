package main

import (
	"flag"
	"net/http"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/providers"
	"github.com/MichaelPalmer1/simple-api-go/providers/base"
	"github.com/MichaelPalmer1/simple-api-go/providers/gcp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func init() {
	api = gcp.FirestoreAPI{
		SimpleAPI: &base.SimpleAPI{
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
	flag.StringVar(&api.Config.OIDCUsernameClaim, "oidc-username-claim", "Sub", "Username claim from OIDC")
	flag.StringVar(&api.Config.OIDCNameClaim, "oidc-name-claim", "Name", "Name claim from OIDC")
	flag.StringVar(&api.Config.OIDCEmailClaim, "oidc-email-claim", "Mail", "Email claim from OIDC")
	flag.StringVar(&api.Config.OIDCGroupClaim, "oidc-group-claim", "", "Group claim from OIDC")
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
	api.Init("simple-api-265401", option.WithCredentialsFile("/home/michael/Downloads/gcp.json"))
	defer api.Close()

	// Initialize http server
	router, err := providers.InitHTTPServer(api, "id", "/items/", []string{"CREATE", "UPDATE"})
	if err != nil {
		panic(err)
	}

	// Add get/create/update endpoints
	router.POST("/item/", create)
	router.GET("/item/:id", get)
	router.PUT("/item/:id", update)
	router.GET("/types/", listTypes)

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", router))
}

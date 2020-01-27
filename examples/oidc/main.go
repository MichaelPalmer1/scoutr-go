package main

import (
	"flag"
	"net/http"
	"os/user"
	"path/filepath"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/providers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

// Initialize - Creates connection to DynamoDB
func Initialize(config *config.Config) *dynamodb.DynamoDB {
	usr, _ := user.Current()

	creds := credentials.NewSharedCredentials(filepath.Join(usr.HomeDir, ".aws/credentials"), "default")
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: creds,
	}))

	svc := dynamodb.New(sess)

	return svc
}

func main() {
	// Command line arguments
	var config config.Config
	flag.StringVar(&config.DataTable, "data-table", "", "Data table")
	flag.StringVar(&config.AuthTable, "auth-table", "", "Auth table")
	flag.StringVar(&config.GroupTable, "group-table", "", "Group table")
	flag.StringVar(&config.AuditTable, "audit-table", "", "Audit table")
	flag.IntVar(&config.LogRetentionDays, "log-retention-days", 30, "Days to retain read logs")
	flag.StringVar(&config.OIDCUsernameClaim, "oidc-username-claim", "Sub", "Username claim from OIDC")
	flag.StringVar(&config.OIDCNameClaim, "oidc-name-claim", "Name", "Name claim from OIDC")
	flag.StringVar(&config.OIDCEmailClaim, "oidc-email-claim", "Mail", "Email claim from OIDC")
	flag.StringVar(&config.OIDCGroupClaim, "oidc-group-claim", "", "Group claim from OIDC")
	flag.Parse()

	// Make sure required fields are provided
	if config.DataTable == "" {
		log.Fatalln("data-table argument is required")
	}
	if config.AuthTable == "" {
		log.Fatalln("auth-table argument is required")
	}
	if config.GroupTable == "" {
		log.Fatalln("group-table argument is required")
	}

	svc := Initialize(&config)
	api.Client = svc
	api.Config = config

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

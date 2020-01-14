package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/endpoints"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Record : Item in Dynamo
type Record map[string]interface{}

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
	flag.Parse()

	svc := Initialize(&config)

	requestUser := models.RequestUser{
		ID: "michael",
	}

	request := models.Request{
		User: requestUser,
	}

	pathParams := make(map[string]string)
	queryParams := make(map[string]string)

	api := endpoints.SimpleAPI{
		DataTable: config.DataTable,
		Client: svc,
	}

	data := api.ListTable(request, "", pathParams, queryParams)

	out, _ := json.Marshal(data)
	fmt.Println(string(out))
}

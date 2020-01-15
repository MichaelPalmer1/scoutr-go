package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
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

var api endpoints.SimpleAPI

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

func httpHandler(w http.ResponseWriter, req *http.Request) {
	pathParams := make(map[string]string)
	queryParams := make(map[string]string)

	requestUser := models.RequestUser{
		ID: "michael",
	}

	// Build the request model
	request := models.Request{
		User:   requestUser,
		Method: req.Method,
		Path:   req.URL.Path,
	}

	if req.Method == "GET" {
		// Parse query params for GET requests
		for key, values := range req.URL.Query() {
			queryParams[key] = values[0]
		}
	} else if req.Method == "POST" || req.Method == "PUT" {
		// Parse the request body if this is a POST/PUT
		var body map[string]interface{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			if err.Error() == "EOF" {
				http.Error(w, "Missing request body", http.StatusBadRequest)
				return
			}

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		request.Body = body
	}

	// List the table
	data, err := api.ListTable(request, "", pathParams, queryParams)

	// Check for errors in the response
	if err != nil {
		switch err.(type) {
		case *models.Unauthorized:
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case *models.BadRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Marshal the response and write it to output
	out, _ := json.Marshal(data)
	w.Header().Add("Content-Type", "application/json")
	w.Write(out)
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
	api.DataTable = config.DataTable
	api.Client = svc

	http.HandleFunc("/", httpHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

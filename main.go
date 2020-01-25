package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os/user"
	"path/filepath"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/providers"
	"github.com/MichaelPalmer1/simple-api-go/simpleapi"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// Record : Item in Dynamo
type Record map[string]interface{}

var api simpleapi.SimpleAPI
var validation map[string]utils.FieldValidation

func init() {
	validation = map[string]utils.FieldValidation{
		"value": func(value string, item map[string]string, existingItem map[string]string) (bool, string, error) {
			if value != "hello" {
				return false, fmt.Sprintf("Invalid value '%s' for attribute 'value'", value), nil
			}

			return true, "", nil
		},
	}
}

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

func create(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	requestUser := models.RequestUser{
		ID: "michael",
	}

	// Parse the request body
	var body map[string]string
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Build the request model
	request := models.Request{
		User:      requestUser,
		Method:    req.Method,
		Path:      req.URL.Path,
		Body:      body,
		SourceIP:  req.RemoteAddr,
		UserAgent: req.UserAgent(),
	}

	// Create the item
	err = api.Create(request, body, validation)

	// Check for errors in the response
	if providers.HTTPErrorHandler(err, w) {
		return
	}

	// Marshal the response and write it to output
	out, _ := json.Marshal(map[string]bool{
		"created": true,
	})
	w.Header().Add("Content-Type", "application/json")
	w.Write(out)
}

func get(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	requestUser := models.RequestUser{
		ID: "michael",
	}

	// Build the request model
	request := models.Request{
		User:      requestUser,
		Method:    req.Method,
		Path:      req.URL.Path,
		SourceIP:  req.RemoteAddr,
		UserAgent: req.UserAgent(),
	}

	// Fetch the item
	data, err := api.Get(request, params.ByName("id"))

	// Check for errors in the response
	if providers.HTTPErrorHandler(err, w) {
		return
	}

	// Marshal the response and write it to output
	out, _ := json.Marshal(data)
	w.Header().Add("Content-Type", "application/json")
	w.Write(out)
}

func update(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	requestUser := models.RequestUser{
		ID: "michael",
	}

	// Parse the request body
	var body map[string]string
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Build the request model
	request := models.Request{
		User:      requestUser,
		Method:    req.Method,
		Path:      req.URL.Path,
		Body:      body,
		SourceIP:  req.RemoteAddr,
		UserAgent: req.UserAgent(),
	}

	// Get key schema
	tableInfo, err := api.Client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(api.Config.DataTable),
	})
	if err != nil {
		log.Errorln("Failed to describe table", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build partition key
	partitionKey := make(map[string]string)
	for _, schema := range tableInfo.Table.KeySchema {
		if *schema.KeyType == "HASH" {
			partitionKey[*schema.AttributeName] = params.ByName("id")
			break
		}
	}

	// Update the item
	data, err := api.Update(request, partitionKey, body, validation, "UPDATE")

	// Check for errors in the response
	if providers.HTTPErrorHandler(err, w) {
		return
	}

	// Marshal the response and write it to output
	out, _ := json.Marshal(data)
	w.Header().Add("Content-Type", "application/json")
	w.Write(out)
}

func listTypes(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	pathParams := make(map[string]string)
	queryParams := make(map[string]string)

	userData := models.UserData{
		Name:     "Michael",
		Email:    "Michael@Palmer.com",
		Username: "michael",
		Groups:   []string{"group1", "group2"},
	}

	requestUser := models.RequestUser{
		ID:   "michael",
		Data: &userData,
	}

	// Parse query params
	for key, values := range req.URL.Query() {
		queryParams[key] = values[0]
	}

	// Parse path params
	for _, item := range params {
		pathParams[item.Key] = item.Value
	}

	// Build the request model
	request := models.Request{
		User:        requestUser,
		Method:      req.Method,
		Path:        req.URL.Path,
		PathParams:  pathParams,
		QueryParams: queryParams,
		SourceIP:    req.RemoteAddr,
		UserAgent:   req.UserAgent(),
	}

	// List the table
	data, err := api.ListUniqueValues(request, "type")

	// Check for errors in the response
	if providers.HTTPErrorHandler(err, w) {
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

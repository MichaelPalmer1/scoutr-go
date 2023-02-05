package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MichaelPalmer1/scoutr-go/pkg/helpers"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func update(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	requestUser := helpers.GetUserFromOIDC(req, api)

	// Parse the request body
	var body map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Parse path params
	pathParams := map[string]string{}
	for _, item := range params {
		pathParams[item.Key] = item.Value
	}

	// Build the request model
	request := types.Request{
		User:       requestUser,
		Method:     req.Method,
		Path:       req.URL.Path,
		PathParams: pathParams,
		Body:       body,
		SourceIP:   req.RemoteAddr,
		UserAgent:  req.UserAgent(),
	}

	// Get key schema
	tableInfo, err := api.Client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(api.Config.DataTable),
	})
	if err != nil {
		log.Errorln("Failed to describe table", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build partition key
	partitionKey := make(map[string]interface{})
	for _, schema := range tableInfo.Table.KeySchema {
		if schema.KeyType == dynamoTypes.KeyTypeHash {
			partitionKey[*schema.AttributeName] = params.ByName("id")
			break
		}
	}

	// Update the item
	data, err := api.Update(request, partitionKey, body, validation, nil, "UPDATE")

	// Check for errors in the response
	if helpers.HTTPErrorHandler(err, w) {
		return
	}

	// Marshal the response and write it to output
	out, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

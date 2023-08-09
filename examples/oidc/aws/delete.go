package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MichaelPalmer1/scoutr-go/helpers"
	"github.com/MichaelPalmer1/scoutr-go/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func delete(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	requestUser := helpers.GetUserFromOIDC(req, api)

	// Build the request model
	request := models.Request{
		User:      requestUser,
		Method:    req.Method,
		Path:      req.URL.Path,
		SourceIP:  req.RemoteAddr,
		UserAgent: req.UserAgent(),
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
	partitionKey := make(map[string]string)
	for _, schema := range tableInfo.Table.KeySchema {
		if schema.KeyType == types.KeyTypeHash {
			partitionKey[*schema.AttributeName] = params.ByName("id")
			break
		}
	}

	// Delete the item
	err = api.Delete(request, partitionKey)

	// Check for errors in the response
	if helpers.HTTPErrorHandler(err, w) {
		return
	}

	// Marshal the response and write it to output
	out, _ := json.Marshal(true)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

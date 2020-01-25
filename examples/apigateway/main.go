package main

import (
	"os"
	"strconv"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/providers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func list(event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// Convert log retention to int
	logRetention, err := strconv.Atoi(os.Getenv("LogRetentionDays"))
	if err != nil {
		panic(err)
	}

	// Build config
	config := config.Config{
		DataTable:        os.Getenv("DataTable"),
		AuthTable:        os.Getenv("AuthTable"),
		AuditTable:       os.Getenv("AuditTable"),
		GroupTable:       os.Getenv("GroupTable"),
		LogRetentionDays: logRetention,
	}

	// Initialize api gateway
	api, request := providers.InitAPIGateway(event, config)

	// List the data
	data, err := api.ListTable(request, "", event.PathParameters, event.QueryStringParameters)

	// Handle any errors
	if errorResponse := providers.APIGatewayErrorHandler(err); errorResponse != nil {
		return *errorResponse
	}

	// Send response
	return providers.ProcessAPIGatewayResponse(data)
}

func get(event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// Convert log retention to int
	logRetention, err := strconv.Atoi(os.Getenv("LogRetentionDays"))
	if err != nil {
		panic(err)
	}

	// Build config
	config := config.Config{
		DataTable:        os.Getenv("DataTable"),
		AuthTable:        os.Getenv("AuthTable"),
		AuditTable:       os.Getenv("AuditTable"),
		GroupTable:       os.Getenv("GroupTable"),
		LogRetentionDays: logRetention,
	}

	// Initialize api gateway
	api, request := providers.InitAPIGateway(event, config)

	// Get the record
	data, err := api.Get(request, event.PathParameters["id"])

	// Handle any errors
	if errorResponse := providers.APIGatewayErrorHandler(err); errorResponse != nil {
		return *errorResponse
	}

	// Send response
	return providers.ProcessAPIGatewayResponse(data)
}

func main() {
	lambda.Start(list)
}

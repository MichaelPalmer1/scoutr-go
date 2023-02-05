package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/MichaelPalmer1/scoutr-go/pkg/config"
	"github.com/MichaelPalmer1/scoutr-go/pkg/helpers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func list(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
	api, request := helpers.InitAPIGateway(event, config)

	// List the data
	data, err := api.List(request)

	// Handle any errors
	if errorResponse := helpers.APIGatewayErrorHandler(err); errorResponse != nil {
		return *errorResponse, nil
	}

	// Send response
	return helpers.ProcessAPIGatewayResponse(data)
}

func local() {
	identity := events.APIGatewayRequestIdentity{
		APIKeyID:  "qtz1hy6j23",
		SourceIP:  "1.2.3.4",
		UserAgent: "Fake",
	}

	context := events.APIGatewayProxyRequestContext{
		Identity: identity,
	}

	event := events.APIGatewayProxyRequest{
		Path:           "/items",
		HTTPMethod:     "GET",
		RequestContext: context,
	}

	resp, err := list(event)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func main() {
	// local()
	lambda.Start(list)
}

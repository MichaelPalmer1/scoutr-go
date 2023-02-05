package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MichaelPalmer1/scoutr-go/pkg/config"
	dynamo "github.com/MichaelPalmer1/scoutr-go/pkg/providers/aws"
	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	log "github.com/sirupsen/logrus"
)

// InitAPIGateway : Initialize API Gateway
func InitAPIGateway(event events.APIGatewayProxyRequest, config config.Config) (dynamo.DynamoAPI, types.Request) {
	// Build request user
	requestUser := types.RequestUser{
		ID: event.RequestContext.Identity.APIKeyID,
	}

	// Build request
	request := types.Request{
		Method:      event.HTTPMethod,
		Path:        event.Path,
		Body:        event.Body,
		PathParams:  event.PathParameters,
		QueryParams: event.MultiValueQueryStringParameters,
		UserAgent:   event.RequestContext.Identity.UserAgent,
		SourceIP:    event.RequestContext.Identity.SourceIP,
		User:        requestUser,
	}

	// Make sure maps are initialized
	if len(event.PathParameters) == 0 {
		request.PathParams = make(map[string]string)
	}
	if len(event.QueryStringParameters) == 0 {
		request.QueryParams = make(map[string][]string)
	}

	// Create API
	api := dynamo.DynamoAPI{
		Scoutr: &base.Scoutr{
			Config: config,
		},
	}

	api.Init(*aws.NewConfig())

	return api, request
}

// APIGatewayErrorHandler : Handle api gateway errors
func APIGatewayErrorHandler(err error) *events.APIGatewayProxyResponse {
	response := &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}

	if err != nil {
		switch err.(type) {
		case *types.Unauthorized:
			response.StatusCode = http.StatusUnauthorized
		case *types.BadRequest:
			response.StatusCode = http.StatusBadRequest
		case *types.NotFound:
			response.StatusCode = http.StatusNotFound
		default:
			response.StatusCode = http.StatusInternalServerError
		}
		log.Errorln("Encountered error", err)
		response.Body = fmt.Sprintf("%s", err)
		return response
	}
	return nil
}

// ProcessAPIGatewayResponse : Process response for api gateway
func ProcessAPIGatewayResponse(data interface{}) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{},
	}
	body, err := json.Marshal(data)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = "Failed to marshal output"
		response.Headers["Content-Type"] = "text/plain"
		log.Errorln("Error marshalling output", err)
		return response, nil
	}

	// Add header and body
	response.Headers["Content-Type"] = "application/json"
	response.Body = string(body)
	response.StatusCode = http.StatusOK

	return response, nil
}

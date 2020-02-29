package helpers

import (
	"net/http"

	"github.com/MichaelPalmer1/simple-api-go/models"
)

// TODO: Add helper functions for Google Cloud Endpoints

// InitCloudEndpoint : Initialize a Google Cloud Endpoint
func InitCloudEndpoint(r *http.Request) models.Request {
	// Build request user
	requestUser := models.RequestUser{
		// ID: event.RequestContext.Identity.APIKeyID,
	}

	// Parse query params
	queryParams := make(map[string]string)
	for key, values := range r.URL.Query() {
		queryParams[key] = values[0]
	}

	// Create the request
	request := models.Request{
		Method: r.Method,
		Path:   r.URL.Path,
		Body:   r.Body,
		// PathParams:  event.PathParameters,
		QueryParams: queryParams,
		UserAgent:   r.UserAgent(),
		SourceIP:    r.RemoteAddr,
		User:        requestUser,
	}

	return request
}

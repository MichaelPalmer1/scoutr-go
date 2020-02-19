package main

import (
	"encoding/json"
	"net/http"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/providers"
	"github.com/julienschmidt/httprouter"
)

func listTypes(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	pathParams := make(map[string]string)
	queryParams := make(map[string]string)

	requestUser := providers.GetUserFromOIDC(req, api)

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
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

package main

import (
	"encoding/json"
	"net/http"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/helpers"
	"github.com/julienschmidt/httprouter"
)

func update(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	requestUser := helpers.GetUserFromOIDC(req, api)

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

	// Build partition key
	partitionKey := map[string]string{
		api.GetConfig().PrimaryKey: params.ByName("id"),
	}

	// Update the item
	data, err := api.Update(request, partitionKey, body, validation, "UPDATE")

	// Check for errors in the response
	if helpers.HTTPErrorHandler(err, w) {
		return
	}

	// Marshal the response and write it to output
	out, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

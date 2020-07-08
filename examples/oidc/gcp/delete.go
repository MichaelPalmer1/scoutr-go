package main

import (
	"encoding/json"
	"net/http"

	"github.com/MichaelPalmer1/scoutr-go/helpers"
	"github.com/MichaelPalmer1/scoutr-go/models"
	"github.com/julienschmidt/httprouter"
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

	// Build partition key
	partitionKey := map[string]string{
		api.Config.PrimaryKey: params.ByName("id"),
	}

	// Delete the item
	err := api.Delete(request, partitionKey)

	// Check for errors in the response
	if helpers.HTTPErrorHandler(err, w) {
		return
	}

	// Marshal the response and write it to output
	out, _ := json.Marshal(true)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

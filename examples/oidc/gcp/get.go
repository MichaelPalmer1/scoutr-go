package main

import (
	"encoding/json"
	"net/http"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/providers"
	"github.com/julienschmidt/httprouter"
)

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
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

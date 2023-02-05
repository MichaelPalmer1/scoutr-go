package main

import (
	"encoding/json"
	"net/http"

	"github.com/MichaelPalmer1/scoutr-go/pkg/helpers"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/julienschmidt/httprouter"
)

func get(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	requestUser := helpers.GetUserFromOIDC(req, api)

	// Build the request model
	request := types.Request{
		User:      requestUser,
		Method:    req.Method,
		Path:      req.URL.Path,
		SourceIP:  req.RemoteAddr,
		UserAgent: req.UserAgent(),
	}

	// Fetch the item
	data, err := api.Get(request, params.ByName("id"))

	// Check for errors in the response
	if helpers.HTTPErrorHandler(err, w) {
		return
	}

	// Marshal the response and write it to output
	out, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

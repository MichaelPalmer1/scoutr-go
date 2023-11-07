package main

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/MichaelPalmer1/scoutr-go/pkg/helpers"
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
// 	"github.com/julienschmidt/httprouter"
// )

// func listTypes(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
// 	pathParams := make(map[string]string)
// 	queryParams := make(map[string][]string)

// 	requestUser := helpers.GetUserFromOIDC(req, api)

// 	// Parse query params
// 	for key, values := range req.URL.Query() {
// 		queryParams[key] = values
// 	}

// 	// Parse path params
// 	for _, item := range params {
// 		pathParams[item.Key] = item.Value
// 	}

// 	// Build the request model
// 	request := types.Request{
// 		User:        requestUser,
// 		Method:      req.Method,
// 		Path:        req.URL.Path,
// 		PathParams:  pathParams,
// 		QueryParams: queryParams,
// 		SourceIP:    req.RemoteAddr,
// 		UserAgent:   req.UserAgent(),
// 	}

// 	// List the table
// 	data, err := api.ListUniqueValues(request, "type")

// 	// Check for errors in the response
// 	if helpers.HTTPErrorHandler(err, w) {
// 		return
// 	}

// 	// Marshal the response and write it to output
// 	out, _ := json.Marshal(data)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(out)
// }

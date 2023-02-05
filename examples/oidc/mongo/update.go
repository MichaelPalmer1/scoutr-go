package main

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/MichaelPalmer1/scoutr-go/pkg/helpers"
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
// 	"github.com/julienschmidt/httprouter"
// )

// func update(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
// 	requestUser := helpers.GetUserFromOIDC(req, api)

// 	// Parse the request body
// 	var body map[string]string
// 	err := json.NewDecoder(req.Body).Decode(&body)
// 	if err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	// Parse path params
// 	pathParams := map[string]string{}
// 	for _, item := range params {
// 		pathParams[item.Key] = item.Value
// 	}

// 	// Build the request model
// 	request := types.Request{
// 		User:       requestUser,
// 		Method:     req.Method,
// 		Path:       req.URL.Path,
// 		PathParams: pathParams,
// 		Body:       body,
// 		SourceIP:   req.RemoteAddr,
// 		UserAgent:  req.UserAgent(),
// 	}

// 	// Build partition key
// 	partitionKey := map[string]string{
// 		api.Config.PrimaryKey: params.ByName("id"),
// 	}

// 	// Update the item
// 	data, err := api.Update(request, partitionKey, body, validation, "UPDATE")

// 	// Check for errors in the response
// 	if helpers.HTTPErrorHandler(err, w) {
// 		return
// 	}

// 	// Marshal the response and write it to output
// 	out, _ := json.Marshal(data)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(out)
// }

package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type userAccess struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// HTTPErrorHandler : Handle HTTP errors
func HTTPErrorHandler(err error, w http.ResponseWriter) bool {
	if err != nil {
		// Marshal the error
		var errorString string
		bs, err := json.Marshal(err)
		if err != nil {
			log.WithError(err).Error("Failed to marshal output")
			errorString = err.Error()
		} else {
			errorString = string(bs)
		}

		// Select the error code
		var errorCode int
		switch err.(type) {
		case *types.Unauthorized:
			errorCode = http.StatusUnauthorized
		case *types.Forbidden:
			errorCode = http.StatusForbidden
		case *types.BadRequest:
			errorCode = http.StatusBadRequest
		case *types.NotFound:
			errorCode = http.StatusNotFound
		default:
			errorCode = http.StatusInternalServerError
		}

		// Trigger the error
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(errorCode)
		if _, err := fmt.Fprintln(w, errorString); err != nil {
			log.WithError(err).Error("Failed to write error content")
		}

		return true
	}
	return false
}

func BuildHttpRequest(api base.ScoutrBase, r *http.Request, params httprouter.Params) types.Request {
	pathParams := make(map[string]string)
	queryParams := make(map[string][]string)

	// Parse query params
	for key, values := range r.URL.Query() {
		queryParams[key] = values
	}

	// Parse path params
	for _, item := range params {
		pathParams[item.Key] = item.Value
	}

	// Parse the request body
	var body interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		body = nil
	}

	// Build request
	req := types.Request{
		User:        GetUserFromOIDC(r, api),
		Method:      r.Method,
		Path:        r.URL.Path,
		Body:        body,
		SourceIP:    r.RemoteAddr,
		UserAgent:   r.UserAgent(),
		PathParams:  pathParams,
		QueryParams: queryParams,
	}

	return req
}

// InitHTTPServer : Initialize the HTTP server
func InitHTTPServer(api base.ScoutrBase, primaryListEndpoint string) (*httprouter.Router, error) {
	// Format primary endpoint
	if !strings.HasPrefix(primaryListEndpoint, "/") {
		primaryListEndpoint = "/" + primaryListEndpoint
	}
	if !strings.HasSuffix(primaryListEndpoint, "/") {
		primaryListEndpoint += "/"
	}

	if strings.Contains(primaryListEndpoint, ":") {
		return nil, errors.New("Path arguments not permitted in primary endpoint")
	}

	list := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// Build request
		request := BuildHttpRequest(api, req, params)

		// List the table
		data, err := api.List(request)

		// Check for errors in the response
		if HTTPErrorHandler(err, w) {
			return
		}

		// Marshal the response and write it to output
		out, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(out)
		if err != nil {
			log.Errorf("Error writing output: %v", err)
		}
	}

	search := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// Build request
		request := BuildHttpRequest(api, req, params)
		values, ok := request.Body.([]string)
		if !ok {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Search the table
		data, err := api.Search(request, params.ByName("key"), values)

		// Check for errors in the response
		if HTTPErrorHandler(err, w) {
			return
		}

		// Marshal the response and write it to output
		out, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(out)
		if err != nil {
			log.Errorf("Error writing output: %v", err)
		}
	}

	userInfo := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// Lookup information about the user
		user := GetUserFromOIDC(req, api)

		// Marshal data and write to output
		data, err := json.Marshal(map[string]string{
			"id":    user.ID,
			"name":  user.Data.Name,
			"email": user.Data.Email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set content type
		w.Header().Set("Content-Type", "application/json")

		// Write output
		_, err = w.Write(data)
		if err != nil {
			log.Errorf("Error writing output: %v", err)
		}
	}

	userHasPermission := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// Build request
		request := BuildHttpRequest(api, req, params)

		// Parse the request body
		access, ok := request.Body.(userAccess)
		if !ok {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Check for authorization and save to output object
		output := map[string]bool{
			"authorized": api.CanAccessEndpoint(access.Method, access.Path, nil, &request),
		}

		// Marshal data and write to output
		data, err := json.Marshal(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set content type
		w.Header().Set("Content-Type", "application/json")

		// Write output
		_, err = w.Write(data)
		if err != nil {
			log.Errorf("Error writing output: %v", err)
		}
	}

	audit := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// Build request
		request := BuildHttpRequest(api, req, params)

		// List the table
		data, err := api.ListAuditLogs(request, request.PathParams, request.QueryParams)

		// Check for errors in the response
		if HTTPErrorHandler(err, w) {
			return
		}

		// Marshal the response and write it to output
		out, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(out)
		if err != nil {
			log.Errorf("Error writing output: %v", err)
		}
	}

	history := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// Build request
		request := BuildHttpRequest(api, req, params)

		// List the table
		data, err := api.History(request, "id", params.ByName("item"), request.QueryParams, []string{"CREATE", "UPDATE", "DELETE"})

		// Check for errors in the response
		if HTTPErrorHandler(err, w) {
			return
		}

		// Marshal the response and write it to output
		out, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(out)
		if err != nil {
			log.Errorf("Error writing output: %v", err)
		}
	}

	// Create routes
	router := httprouter.New()
	router.GET(primaryListEndpoint, list)
	router.GET(primaryListEndpoint+":search_key/:search_value/", list)
	router.GET("/user/", userInfo)
	router.POST("/user/has-permission/", userHasPermission)
	router.GET("/audit/", audit)
	router.GET("/audit/:item/", audit)
	router.GET("/history/:item/", history)
	router.POST("/search/:key/", search)

	return router, nil
}

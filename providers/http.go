package providers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/simpleapi"
	"github.com/julienschmidt/httprouter"
)

type userAccess struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// HTTPErrorHandler : Handle HTTP errors
func HTTPErrorHandler(err error, w http.ResponseWriter) bool {
	if err != nil {
		switch err.(type) {
		case *models.Unauthorized:
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case *models.BadRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case *models.NotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return true
	}
	return false
}

// InitHTTPServer : Initialize the HTTP server
func InitHTTPServer(api simpleapi.SimpleAPI, partitionKey string, primaryListEndpoint string, historyActions []string) (*httprouter.Router, error) {

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
		pathParams := make(map[string]string)
		queryParams := make(map[string]string)

		userData := models.UserData{
			Name:     "Michael",
			Email:    "Michael@Palmer.com",
			Username: "michael",
			Groups:   []string{"group1", "group2"},
		}

		requestUser := models.RequestUser{
			ID:   "michael",
			Data: &userData,
		}

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
		data, err := api.ListTable(request)

		// Check for errors in the response
		if HTTPErrorHandler(err, w) {
			return
		}

		// Marshal the response and write it to output
		out, _ := json.Marshal(data)
		w.Header().Add("Content-Type", "application/json")
		w.Write(out)
	}

	search := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		requestUser := models.RequestUser{
			ID: "michael",
		}

		// Parse the request body
		var body []string
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

		// Search the table
		data, err := api.Search(request, params.ByName("key"), body)

		// Check for errors in the response
		if HTTPErrorHandler(err, w) {
			return
		}

		// Marshal the response and write it to output
		out, _ := json.Marshal(data)
		w.Header().Add("Content-Type", "application/json")
		w.Write(out)
	}

	userInfo := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// Lookup information about the user
		info := map[string]string{
			"id":    req.Header.Get("Oidc-Claim-Sub"),
			"name":  req.Header.Get("Oidc-Claim-Name"),
			"email": req.Header.Get("Oidc-Claim-Mail"),
		}

		// Marshal data and write to output
		data, err := json.Marshal(info)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set content type
		w.Header().Set("Content-Type", "application/json")

		// Write output
		w.Write(data)
	}

	userHasPermission := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		userData := models.UserData{
			Name:     req.Header.Get("Oidc-Claim-Name"),
			Email:    req.Header.Get("Oidc-Claim-Mail"),
			Username: req.Header.Get("Oidc-Claim-Sub"),
		}

		requestUser := models.RequestUser{
			ID:   req.Header.Get("Oidc-Claim-Sub"),
			Data: &userData,
		}

		// Build the request model
		request := models.Request{
			User:   requestUser,
			Method: req.Method,
			Path:   req.URL.Path,
		}
		_ = request

		// Parse the request body
		access := userAccess{}
		err := json.NewDecoder(req.Body).Decode(&access)
		if err != nil {
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
		w.Write(data)
	}

	audit := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		pathParams := make(map[string]string)
		queryParams := make(map[string]string)

		userData := models.UserData{
			Name:     "Michael",
			Email:    "Michael@Palmer.com",
			Username: "michael",
			Groups:   []string{"group1", "group2"},
		}

		requestUser := models.RequestUser{
			ID:   "michael",
			Data: &userData,
		}

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
			QueryParams: queryParams,
			SourceIP:    req.RemoteAddr,
			UserAgent:   req.UserAgent(),
		}

		// List the table
		data, err := api.ListAuditLogs(request, pathParams, queryParams)

		// Check for errors in the response
		if HTTPErrorHandler(err, w) {
			return
		}

		// Marshal the response and write it to output
		out, _ := json.Marshal(data)
		w.Header().Add("Content-Type", "application/json")
		w.Write(out)
	}

	history := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		queryParams := make(map[string]string)

		userData := models.UserData{
			Name:     "Michael",
			Email:    "Michael@Palmer.com",
			Username: "michael",
			Groups:   []string{"group1", "group2"},
		}

		requestUser := models.RequestUser{
			ID:   "michael",
			Data: &userData,
		}

		// Parse query params
		for key, values := range req.URL.Query() {
			queryParams[key] = values[0]
		}

		// Build the request model
		request := models.Request{
			User:        requestUser,
			Method:      req.Method,
			Path:        req.URL.Path,
			QueryParams: queryParams,
			SourceIP:    req.RemoteAddr,
			UserAgent:   req.UserAgent(),
		}

		// List the table
		data, err := api.History(request, "id", params.ByName("item"), queryParams, []string{"CREATE", "UPDATE", "DELETE"})

		// Check for errors in the response
		if HTTPErrorHandler(err, w) {
			return
		}

		// Marshal the response and write it to output
		out, _ := json.Marshal(data)
		w.Header().Add("Content-Type", "application/json")
		w.Write(out)
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

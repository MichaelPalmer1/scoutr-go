package azure

import (
	"time"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

// ListAuditLogs : List audit logs
func (api MongoDBAPI) ListAuditLogs(req models.Request, pathParams map[string]string, queryParams map[string]string) ([]models.AuditLog, error) {
	// Only fetch audit logs if the table is configured
	if api.Config.AuditTable == "" {
		return nil, &models.NotFound{
			Message: "Audit logs are not enabled",
		}
	}

	// Get the user
	_, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Generate dynamic search
	searchKey, hasSearchKey := pathParams["search_key"]
	searchValue, hasSearchValue := pathParams["search_value"]
	if hasSearchKey && hasSearchValue {
		// Map the search key and value into path params
		pathParams[searchKey] = searchValue
		delete(pathParams, "search_key")
		delete(pathParams, "search_value")
	}

	// Merge pathParams into queryParams
	for key, value := range pathParams {
		queryParams[key] = value
	}

	// Create collection and filter
	collection := api.Client.C(api.Config.AuditTable)
	rawFilters, _, err := api.Filter(&api.Filtering, nil, queryParams)
	if err != nil {
		return nil, err
	}

	// Make sure filters are cast as bson.D
	var filters bson.D
	if rawFilters == nil {
		// No filters
		filters = nil
	} else if _, ok := rawFilters.(bson.DocElem); ok {
		// Single filter
		filters = bson.D{rawFilters.(bson.DocElem)}
	} else {
		// Multiple filters
		filters = rawFilters.(bson.D)
	}

	// Query the data
	var records []models.AuditLog
	err = collection.Find(filters).Sort("time").All(&records)
	if err != nil {
		log.Errorf("Error querying collection: %v", err)
		return nil, err
	}

	return records, nil
}

// auditLog : Creates an audit log
func (api MongoDBAPI) auditLog(action string, request models.Request, user models.User, resource *map[string]string, changes *map[string]string) error {
	// Only send audit logs if the table is configured
	if api.Config.AuditTable == "" {
		return nil
	}

	// Create audit log
	now := time.Now().UTC()
	auditLog := models.AuditLog{
		Time: now.Format(time.RFC3339Nano),
		User: models.AuditUser{
			ID:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			SourceIP:  request.SourceIP,
			UserAgent: request.UserAgent,
		},
		Action: action,
		Method: request.Method,
		Path:   request.Path,
	}

	// Add expiry time for read events
	if action == "GET" || action == "LIST" || action == "SEARCH" {
		auditLog.ExpireTime = now.AddDate(0, 0, api.Config.LogRetentionDays).Unix()
	}

	// Add query params
	if len(request.QueryParams) > 0 {
		auditLog.QueryParams = request.QueryParams
	}

	// Add body
	if request.Body != nil {
		auditLog.Body = request.Body
	} else if changes != nil {
		auditLog.Body = *changes
	}

	// Add resource
	if resource != nil {
		auditLog.Resource = *resource
	}

	// Create collection
	collection := api.Client.C(api.Config.AuditTable)

	// Add the record
	err := collection.Insert(auditLog)
	if err != nil {
		log.Errorln("Failed to save audit log", err)
		log.Infof("Failed audit log: '%s'", auditLog)
		return err
	}

	return nil
}

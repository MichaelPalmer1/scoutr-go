package gcp

import (
	"time"

	"cloud.google.com/go/firestore"
	"github.com/MichaelPalmer1/simple-api-go/models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListAuditLogs : List audit logs
func (api FirestoreAPI) ListAuditLogs(req models.Request, pathParams map[string]string, queryParams map[string]string) ([]models.AuditLog, error) {
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
	collection := api.Client.Collection(api.Config.AuditTable)
	f := FirestoreFiltering{
		Query: collection.Query,
	}
	filters, _, err := api.Filter(&f, nil, queryParams)
	if err != nil {
		return nil, err
	}
	query := collection.Query
	if filters != nil {
		query = filters.(firestore.Query)
	}

	// Order by time
	query = query.OrderBy("time", firestore.Desc)

	// Query the data
	iter := query.Documents(api.context)
	records := []models.AuditLog{}

	// Iterate through the results
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			// Attempt to convert error to a status code
			code, ok := status.FromError(err)

			// Check if the status conversion was successful
			if ok {
				switch code.Code() {
				case codes.InvalidArgument:
					// Return bad request on invalid argument errors
					return nil, &models.BadRequest{
						Message: code.Message(),
					}
				}
			}

			// Fallback to just returning the raw error
			return nil, err
		}

		// Cast result to AuditLog
		var auditLog models.AuditLog
		if err = doc.DataTo(&auditLog); err != nil {
			return nil, err
		}

		// Add audit log to output
		records = append(records, auditLog)
	}

	return records, nil
}

// auditLog : Creates an audit log
func (api FirestoreAPI) auditLog(action string, request models.Request, user models.User, resource *map[string]string, changes *map[string]string) error {
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
	collection := api.Client.Collection(api.Config.AuditTable)

	// Add the record
	_, err := collection.Doc(auditLog.Time).Create(api.context, auditLog)
	if err != nil {
		log.Errorln("Failed to save audit log", err)
		log.Infof("Failed audit log: '%v'", auditLog)
		return err
	}

	return nil
}

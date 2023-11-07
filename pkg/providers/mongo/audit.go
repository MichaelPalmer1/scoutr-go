package mongo

// import (
// 	"time"

// 	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
// 	"github.com/globalsign/mgo/bson"
// 	log "github.com/sirupsen/logrus"
// )

// // ListAuditLogs : List audit logs
// func (api MongoAPI) ListAuditLogs(req types.Request, pathParams map[string]string, queryParams map[string][]string) ([]types.AuditLog, error) {
// 	// Only fetch audit logs if the table is configured
// 	if api.Config.AuditTable == "" {
// 		return nil, &types.NotFound{
// 			Message: "Audit logs are not enabled",
// 		}
// 	}

// 	// Get the user
// 	_, err := api.InitializeRequest(req)
// 	if err != nil {
// 		// Bad user - pass the error through
// 		return nil, err
// 	}

// 	// Generate dynamic search
// 	searchKey, hasSearchKey := pathParams["search_key"]
// 	searchValue, hasSearchValue := pathParams["search_value"]
// 	if hasSearchKey && hasSearchValue {
// 		// Map the search key and value into path params
// 		pathParams[searchKey] = searchValue
// 		delete(pathParams, "search_key")
// 		delete(pathParams, "search_value")
// 	}

// 	// Merge pathParams into queryParams
// 	for key, value := range pathParams {
// 		queryParams[key] = []string{value}
// 	}

// 	// Create collection and filter
// 	collection := api.Client.C(api.Config.AuditTable)
// 	rawFilters, err := api.Filtering.Filter(nil, queryParams, base.FilterActionRead)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Make sure filters are cast as bson.D
// 	var filters bson.D
// 	if rawFilters == nil {
// 		// No filters
// 		filters = nil
// 	} else if _, ok := rawFilters.(bson.DocElem); ok {
// 		// Single filter
// 		filters = bson.D{rawFilters.(bson.DocElem)}
// 	} else {
// 		// Multiple filters
// 		filters = rawFilters.(bson.D)
// 	}

// 	// Query the data
// 	var records []types.AuditLog
// 	err = collection.Find(filters).Sort("time").All(&records)
// 	if err != nil {
// 		log.Errorf("Error querying collection: %v", err)
// 		return nil, err
// 	}

// 	return records, nil
// }

// // auditLog : Creates an audit log
// func (api MongoAPI) auditLog(action string, request types.Request, user types.User, resource *map[string]interface{}, changes *map[string]string) error {
// 	// Only send audit logs if the table is configured
// 	if api.Config.AuditTable == "" {
// 		return nil
// 	}

// 	// Create audit log
// 	now := time.Now().UTC()
// 	auditLog := types.AuditLog{
// 		Time: now.Format(time.RFC3339Nano),
// 		User: types.AuditUser{
// 			ID:        user.ID,
// 			Name:      user.Name,
// 			Username:  user.Username,
// 			Email:     user.Email,
// 			SourceIP:  request.SourceIP,
// 			UserAgent: request.UserAgent,
// 		},
// 		Action: action,
// 		Method: request.Method,
// 		Path:   request.Path,
// 	}

// 	// Add expiry time for read events
// 	if action == base.AuditActionGet || action == base.AuditActionList || action == base.AuditActionSearch {
// 		auditLog.ExpireTime = now.AddDate(0, 0, api.Config.LogRetentionDays).Unix()
// 	}

// 	// Add query params
// 	if len(request.QueryParams) > 0 {
// 		auditLog.QueryParams = request.QueryParams
// 	}

// 	// Add body
// 	if request.Body != nil {
// 		auditLog.Body = request.Body
// 	} else if changes != nil {
// 		auditLog.Body = *changes
// 	}

// 	// Add resource
// 	if resource != nil {
// 		auditLog.Resource = *resource
// 	}

// 	// Create collection
// 	collection := api.Client.C(api.Config.AuditTable)

// 	// Add the record
// 	err := collection.Insert(auditLog)
// 	if err != nil {
// 		log.Errorln("Failed to save audit log", err)
// 		log.Infof("Failed audit log: '%v'", auditLog)
// 		return err
// 	}

// 	return nil
// }

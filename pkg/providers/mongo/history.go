package mongo

// import (
// 	"errors"
// 	"sort"

// 	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
// 	"github.com/globalsign/mgo/bson"
// 	log "github.com/sirupsen/logrus"
// )

// // History : Generate record history
// func (api MongoAPI) History(req types.Request, key string, value string, queryParams map[string][]string, actions []string) ([]types.History, error) {
// 	history := []types.History{}

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

// 	searchParams := map[string]string{
// 		"resource." + key: value,
// 		// "action__in": actions,
// 	}

// 	// Get the audit logs
// 	data, err := api.ListAuditLogs(req, searchParams, queryParams)
// 	if err != nil {
// 		log.Errorln("Error listing audit logs", err)
// 		return nil, err
// 	}

// 	// No results
// 	if len(data) == 0 {
// 		return history, nil
// 	}

// 	// Sort the results
// 	sort.Slice(data, func(i, j int) bool {
// 		return data[i].Time < data[j].Time
// 	})

// 	// Find original creation record
// 	var currentItem types.History
// 	for _, item := range data {
// 		if item.Action == base.AuditActionCreate {
// 			body := make(map[string]interface{})
// 			for key, value := range item.Body.(bson.M) {
// 				body[key] = value
// 			}
// 			currentItem.Time = item.Time
// 			currentItem.Data = body
// 			break
// 		}
// 	}

// 	if currentItem.Time == "" {
// 		return history, errors.New("Failed to find initial creation record")
// 	}

// 	// Make a copy of the original record
// 	originalItem := types.History{
// 		Time: currentItem.Time,
// 		Data: types.Record{},
// 	}
// 	for key, value := range currentItem.Data {
// 		originalItem.Data[key] = value
// 	}

// 	// Add original record
// 	history = append(history, originalItem)

// 	// Parse data
// 	for _, item := range data {
// 		// Skip create records
// 		if item.Action == base.AuditActionCreate {
// 			continue
// 		} else if item.Action == base.AuditActionDelete {
// 			// Insert at the top
// 			history = append(history, types.History{})
// 			copy(history[1:], history[0:])
// 			history[0] = types.History{Time: item.Time}
// 			continue
// 		} else if item.Action == base.AuditActionGet || item.Action == base.AuditActionSearch {
// 			// Skip read actions
// 			continue
// 		}

// 		// Update item
// 		for key, value := range item.Body.(map[string]interface{}) {
// 			currentItem.Data[key] = value
// 		}

// 		// Make a copy of the current item
// 		newItem := types.History{
// 			Time: item.Time,
// 			Data: types.Record{},
// 		}
// 		for key, value := range currentItem.Data {
// 			newItem.Data[key] = value
// 		}

// 		// Insert at the top
// 		history = append(history, types.History{})
// 		copy(history[1:], history[0:])
// 		history[0] = newItem
// 	}

// 	return history, nil
// }

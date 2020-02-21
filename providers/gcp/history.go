package gcp

import (
	"errors"
	"sort"

	"github.com/MichaelPalmer1/simple-api-go/models"
	log "github.com/sirupsen/logrus"
)

// History : Generate record history
func (api FirestoreAPI) History(req models.Request, key string, value string, queryParams map[string]string, actions []string) ([]models.History, error) {
	var history []models.History

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

	searchParams := map[string]string{
		"Resource." + key: value,
		// "action__in": actions,
	}

	// Get the audit logs
	data, err := api.ListAuditLogs(req, searchParams, queryParams)
	if err != nil {
		log.Errorln("Error listing audit logs", err)
		return nil, err
	}

	// No results
	if len(data) == 0 {
		return history, nil
	}

	// Sort the results
	sort.Slice(data, func(i, j int) bool {
		return data[i].Time < data[j].Time
	})

	// Find original creation record
	var currentItem models.History
	for _, item := range data {
		if item.Action == "CREATE" {
			currentItem.Time = item.Time
			currentItem.Data = item.Body.(map[string]interface{})
			break
		}
	}

	if currentItem.Time == "" {
		return history, errors.New("Failed to find initial creation record")
	}

	// Make a copy of the original record
	originalItem := models.History{
		Time: currentItem.Time,
		Data: models.Record{},
	}
	for key, value := range currentItem.Data {
		originalItem.Data[key] = value
	}

	// Add original record
	history = append(history, originalItem)

	// Parse data
	for _, item := range data {
		// Skip create records
		if item.Action == "CREATE" {
			continue
		} else if item.Action == "DELETE" {
			// Insert at the top
			history = append(history, models.History{})
			copy(history[1:], history[0:])
			history[0] = models.History{Time: item.Time}
			continue
		} else if item.Action == "GET" || item.Action == "SEARCH" {
			// Skip read actions
			continue
		}

		// Update item
		for key, value := range item.Body.(map[string]interface{}) {
			currentItem.Data[key] = value
		}

		// Make a copy of the current item
		newItem := models.History{
			Time: item.Time,
			Data: models.Record{},
		}
		for key, value := range currentItem.Data {
			newItem.Data[key] = value
		}

		// Insert at the top
		history = append(history, models.History{})
		copy(history[1:], history[0:])
		history[0] = newItem
	}

	return history, nil
}

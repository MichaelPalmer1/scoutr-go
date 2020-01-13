package utils

import "github.com/MichaelPalmer1/simple-api-go/models"

// PostProcess : Perform post processing on data
func PostProcess(data []models.Record, user *models.User) []models.Record {
	for _, item := range data {
		for _, key := range user.ExcludeFields {
			if _, ok := item[key]; ok {
				delete(item, key)
			}
		}
	}

	return data
}

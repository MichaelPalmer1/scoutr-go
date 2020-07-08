package utils

import (
	"strings"

	"github.com/MichaelPalmer1/scoutr-go/models"
)

// FieldValidation : Callable
type FieldValidation func(value string, item map[string]string, existingItem map[string]string) (bool, string, error)

// ValidateFields : Perform field validation
func ValidateFields(validation map[string]FieldValidation, item map[string]string, existingItem map[string]string, ignoreFieldPresence bool) error {
	// Check for required fields
	if !ignoreFieldPresence {
		var missingKeys []string
		for key := range validation {
			if _, ok := item[key]; !ok {
				missingKeys = append(missingKeys, key)
			}
		}
		if len(missingKeys) > 0 {
			return &models.BadRequest{
				Message: "Missing required fields: " + strings.Join(missingKeys, ", "),
			}
		}
	}

	for key, fn := range validation {
		if _, ok := item[key]; ok {
			success, message, err := fn(item[key], item, existingItem)
			if err != nil {
				return err
			} else if !success {
				return &models.BadRequest{
					Message: message,
				}
			}
		}
	}
	return nil
}

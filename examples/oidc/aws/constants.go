package main

import (
	"fmt"

	dynamo "github.com/MichaelPalmer1/scoutr-go/providers/aws"
	"github.com/MichaelPalmer1/scoutr-go/utils"
)

var api dynamo.DynamoAPI
var validation map[string]utils.FieldValidation

func init() {
	validation = map[string]utils.FieldValidation{
		"value": func(value string, item map[string]string, existingItem map[string]string) (bool, string, error) {
			if value != "hello" {
				return false, fmt.Sprintf("Invalid value '%s' for attribute 'value'", value), nil
			}

			return true, "", nil
		},
	}
}

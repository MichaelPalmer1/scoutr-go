package main

import (
	"fmt"

	"github.com/MichaelPalmer1/scoutr-go/models"

	dynamo "github.com/MichaelPalmer1/scoutr-go/providers/aws"
	"github.com/MichaelPalmer1/scoutr-go/utils"
)

var api dynamo.DynamoAPI
var validation map[string]models.FieldValidation

func init() {
	options := []string{"a", "b", "c"}

	validation = map[string]models.FieldValidation{
		"value": func(input *models.ValidationInput, ch chan models.ValidationOutput) {
			if input.Value != "hello" {
				ch <- models.ValidationOutput{
					Input:   input,
					Result:  false,
					Message: fmt.Sprintf("Invalid value '%s' for attribute 'value'", input.Value),
				}
				return
			}

			ch <- models.ValidationOutput{
				Input:  input,
				Result: true,
			}
		},
		"name": utils.ValueInArray(options, "letter", ""),
	}
}

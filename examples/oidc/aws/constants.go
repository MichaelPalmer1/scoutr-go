package main

import (
	"fmt"

	dynamo "github.com/MichaelPalmer1/scoutr-go/pkg/providers/aws"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/MichaelPalmer1/scoutr-go/pkg/utils"
)

var api dynamo.DynamoAPI
var validation map[string]types.FieldValidation

func init() {
	options := []string{"a", "b", "c"}

	validation = map[string]types.FieldValidation{
		"value": func(input *types.ValidationInput, ch chan types.ValidationOutput) {
			if input.Value != "hello" {
				ch <- types.ValidationOutput{
					Input:   input,
					Result:  false,
					Message: fmt.Sprintf("Invalid value '%s' for attribute 'value'", input.Value),
				}
				return
			}

			ch <- types.ValidationOutput{
				Input:  input,
				Result: true,
			}
		},
		"name": utils.ValueInArray(options, "letter", ""),
	}
}

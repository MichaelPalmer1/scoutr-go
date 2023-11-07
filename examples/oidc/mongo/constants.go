package main

// import (
// 	"fmt"

// 	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/mongo"
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
// 	"github.com/MichaelPalmer1/scoutr-go/pkg/utils"
// )

// var api mongo.MongoAPI
// var validation map[string]utils.FieldValidation

// func init() {
// 	validation = map[string]utils.FieldValidation{
// 		"value": func(input *types.ValidationInput, ch chan types.ValidationOutput) {
// 			output := types.ValidationOutput{
// 				Input: input,
// 			}

// 			if input.Value == "hello" {
// 				output.Result = true
// 				ch <- output
// 				return
// 			}

// 			output.Result = false
// 			output.Message = fmt.Sprintf("Invalid value '%s' for attribute 'value'", input.Value)

// 			ch <- output
// 		},
// 	}
// }

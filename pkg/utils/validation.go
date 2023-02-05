package utils

import (
	"fmt"

	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
)

// FieldValidation : Callable
// type FieldValidation func(input *types.ValidationInput, ch chan types.ValidationOutput)

// ValidateFields : Perform field validation
//func ValidateFields(validation map[string]FieldValidation, requiredFields []string, item map[string]string, existingItem map[string]string) error {
//	// Check for required fields
//	if len(requiredFields) > 0 {
//		var missingKeys []string
//		for _, key := range requiredFields {
//			if _, ok := item[key]; !ok {
//				missingKeys = append(missingKeys, key)
//			}
//		}
//		if len(missingKeys) > 0 {
//			return &types.BadRequest{
//				Message: "Missing required fields: " + strings.Join(missingKeys, ", "),
//			}
//		}
//	}
//
//	// Create channels and wait group
//	wg := &sync.WaitGroup{}
//	ch := make(chan types.ValidationOutput, len(validation))
//	done := make(chan bool, 1)
//
//	// Trigger validation goroutines
//	for key, fn := range validation {
//		if _, ok := item[key]; ok {
//			input := &types.ValidationInput{
//				Value:        item[key],
//				Item:         item,
//				ExistingItem: existingItem,
//			}
//
//			// Increment wait group and start goroutine
//			wg.Add(1)
//			go fn(input, ch)
//		}
//	}
//
//	// Wait for all validations to finish
//	go func(ch chan bool) {
//		wg.Wait()
//		ch <- true
//	}(done)
//
//	// Receive results
//	for {
//		select {
//		case output := <-ch:
//			if output.Error != nil {
//				// Validation threw an error, return the error
//				return output.Error
//			} else if !output.Result {
//				// Validation failed, return with the error message
//				return &types.BadRequest{
//					Message: output.Message,
//				}
//			}
//
//			// Complete wait group item
//			wg.Done()
//
//		case <-done:
//			// Return when all validations have been processed
//			return nil
//		}
//
//	}
//}

func ValueInArray(validOptions []string, optionName string, customErrorMessage string) func(*types.ValidationInput, chan types.ValidationOutput) {
	if optionName == "" {
		optionName = "option"
	}

	return func(input *types.ValidationInput, ch chan types.ValidationOutput) {
		for _, item := range validOptions {
			if item == input.Value {
				ch <- types.ValidationOutput{Result: true}
				return
			}
		}

		errorMessage := customErrorMessage
		if errorMessage == "" {
			errorMessage = fmt.Sprintf("%s is not a valid %s. Valid options: %+v", input.Value, optionName, validOptions)
		}

		ch <- types.ValidationOutput{
			Result:  false,
			Message: errorMessage,
			Error:   nil,
		}
	}
}

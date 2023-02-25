package utils_test

import (
	"testing"

	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/MichaelPalmer1/scoutr-go/pkg/utils"
)

func TestValueInArray(t *testing.T) {
	options := []string{"1", "2", "3"}
	fn := utils.ValueInArray(options, "option", "error message")

	input := &types.ValidationInput{
		Key:   "key",
		Value: "1",
	}

	// Create output channel
	output := make(chan types.ValidationOutput)
	defer close(output)

	// Trigger function in goroutine
	go fn(input, output)

	// Get result from channel
	result := <-output

	if result.Error != nil {
		t.Error(result.Error)
	}

	if !result.Result {
		t.Error("Result should be true")
	}
}

func TestValueInArrayBadValue(t *testing.T) {
	options := []string{"1", "2", "3"}
	fn := utils.ValueInArray(options, "", "")

	input := &types.ValidationInput{
		Key:   "key",
		Value: "4",
	}

	// Create output channel
	output := make(chan types.ValidationOutput)
	defer close(output)

	// Trigger function in goroutine
	go fn(input, output)

	// Get result from channel
	result := <-output

	if result.Error != nil {
		t.Error(result.Error)
	}

	if result.Result {
		t.Error("Result should be false")
	}

	expected := "4 is not a valid option. Valid options: [1 2 3]"
	if result.Message != expected {
		t.Errorf("Expected message '%s' but got '%s'", expected, result.Message)
	}
}

func TestValueInArrayBadValueCustomOptionName(t *testing.T) {
	options := []string{"1", "2", "3"}
	fn := utils.ValueInArray(options, "val", "")

	input := &types.ValidationInput{
		Key:   "key",
		Value: "4",
	}

	// Create output channel
	output := make(chan types.ValidationOutput)
	defer close(output)

	// Trigger function in goroutine
	go fn(input, output)

	// Get result from channel
	result := <-output

	if result.Error != nil {
		t.Error(result.Error)
	}

	if result.Result {
		t.Error("Result should be false")
	}

	expected := "4 is not a valid val. Valid options: [1 2 3]"
	if result.Message != expected {
		t.Errorf("Expected message '%s' but got '%s'", expected, result.Message)
	}
}

func TestValueInArrayBadValueCustomErrorMessage(t *testing.T) {
	options := []string{"1", "2", "3"}
	fn := utils.ValueInArray(options, "val", "custom error message")

	input := &types.ValidationInput{
		Key:   "key",
		Value: "4",
	}

	// Create output channel
	output := make(chan types.ValidationOutput)
	defer close(output)

	// Trigger function in goroutine
	go fn(input, output)

	// Get result from channel
	result := <-output

	if result.Error != nil {
		t.Error(result.Error)
	}

	if result.Result {
		t.Error("Result should be false")
	}

	expected := "custom error message"
	if result.Message != expected {
		t.Errorf("Expected message '%s' but got '%s'", expected, result.Message)
	}
}

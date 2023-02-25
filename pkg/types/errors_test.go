package types_test

import (
	"testing"

	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
)

func TestUnauthorized(t *testing.T) {
	err := types.Unauthorized{
		Message: "not authorized",
	}

	if err.Error() != "not authorized" {
		t.Errorf("Expected error message 'not authorized' but got '%s'", err.Error())
	}
}

func TestUnauthorizedMessages(t *testing.T) {
	err := types.Unauthorized{
		Messages: map[string]string{
			"error": "message",
		},
	}

	expected := `{"error":"message"}`
	if err.Error() != expected {
		t.Errorf("Expected error message %s' but got '%s'", expected, err.Error())
	}
}

func TestForbidden(t *testing.T) {
	err := types.Forbidden{
		Message: "not authorized",
	}

	if err.Error() != "not authorized" {
		t.Errorf("Expected error message 'not authorized' but got '%s'", err.Error())
	}
}

func TestForbiddenMessages(t *testing.T) {
	err := types.Forbidden{
		Messages: map[string]string{
			"error": "message",
		},
	}

	expected := `{"error":"message"}`
	if err.Error() != expected {
		t.Errorf("Expected error message %s' but got '%s'", expected, err.Error())
	}
}

func TestBadRequest(t *testing.T) {
	err := types.BadRequest{
		Message: "not authorized",
	}

	if err.Error() != "not authorized" {
		t.Errorf("Expected error message 'not authorized' but got '%s'", err.Error())
	}
}

func TestBadRequestMessages(t *testing.T) {
	err := types.BadRequest{
		Messages: map[string]string{
			"error": "message",
		},
	}

	expected := `{"error":"message"}`
	if err.Error() != expected {
		t.Errorf("Expected error message %s' but got '%s'", expected, err.Error())
	}
}

func TestNotFound(t *testing.T) {
	err := types.NotFound{
		Message: "not authorized",
	}

	if err.Error() != "not authorized" {
		t.Errorf("Expected error message 'not authorized' but got '%s'", err.Error())
	}
}

func TestNotFoundMessages(t *testing.T) {
	err := types.NotFound{
		Messages: map[string]string{
			"error": "message",
		},
	}

	expected := `{"error":"message"}`
	if err.Error() != expected {
		t.Errorf("Expected error message %s' but got '%s'", expected, err.Error())
	}
}

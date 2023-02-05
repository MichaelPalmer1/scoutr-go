package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func TestGenerateInFilter(t *testing.T) {
	var values []string
	for i := 0; i < 100; i++ {
		values = append(values, fmt.Sprintf("%d", i))
	}

	t.Logf("Values: %+v", values)

	conditions := generateInFilter("key", values)

	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Expr Names: %+v", expr.Names())
	t.Logf("Expr Values: %+v", expr.Values())

	if len(expr.Values()) != len(values) {
		t.Errorf("Expected %d values, but got %d", len(values), len(expr.Values()))
	}
}

func TestGenerateInFilterNoValue(t *testing.T) {
	conditions := generateInFilter("key", nil)

	if conditions.IsSet() {
		t.Error("Conditions should not be set")
	}
}

func TestGenerateInFilterSingleValue(t *testing.T) {
	conditions := generateInFilter("key", []string{"value1"})

	if !conditions.IsSet() {
		t.Error("Conditions should be set")
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		t.Fatal(err)
	}

	if len(expr.Values()) != 1 {
		t.Errorf("Expected 1 value  but got %d", len(expr.Values()))
	}
}

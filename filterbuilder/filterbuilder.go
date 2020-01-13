package filterbuilder

import (
	"fmt"
	"reflect"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// BuildFilter : Build a filter
func BuildFilter(user *models.User, filters map[string]string) *expression.ConditionBuilder {
	var conditions expression.ConditionBuilder

	// Build user filters
	for idx, item := range user.FilterFields {
		attr := expression.Name(item.Field)
		if value, ok := item.Value.(string); ok {
			// Value is a single string
			condition := attr.Equal(expression.Value(value))
			if idx == 0 {
				conditions = condition
			} else {
				conditions = conditions.And(condition)
			}
		} else if value, ok := item.Value.([]interface{}); ok {
			// Value is a list of strings
			condition := attr.In(expression.Value(value))
			if idx == 0 {
				conditions = condition
			} else {
				conditions = conditions.And(condition)
			}
		} else {
			fmt.Println("what's this", item.Value)
			fmt.Println(reflect.TypeOf(item.Value))
			continue
		}
	}

	// Build specified filters
	// for _, item := range filters {
	// 	item
	// }

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		fmt.Println("error building condition expression", err)
		return nil
	}

	fmt.Println("Values:", expr.Values())
	fmt.Println("Filter:", *expr.Filter())

	return &conditions
}

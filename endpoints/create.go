package endpoints

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MichaelPalmer1/simple-api-go/filterbuilder"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// FieldValidation : Callable
type FieldValidation func(value string, item map[string]string, existingItem map[string]string) (bool, string, error)

func validateFields(validation map[string]FieldValidation, item map[string]string, existingItem map[string]string, ignoreFieldPresence bool) error {
	// Check for required fields
	if !ignoreFieldPresence {
		var missingKeys []string
		for key := range validation {
			if _, ok := item[key]; !ok {
				missingKeys = append(missingKeys, key)
			}
		}
		if len(missingKeys) > 0 {
			return errors.New("Missing required fields: " + strings.Join(missingKeys, ", "))
		}
	}

	for key, fn := range validation {
		if _, ok := item[key]; ok {
			success, message, err := fn(item[key], item, existingItem)
			if err != nil {
				return err
			} else if !success {
				return errors.New(message)
			}
		}
	}
	return nil
}

// Create : Create an item
func (api *SimpleAPI) Create(req models.Request, item map[string]string, validation map[string]FieldValidation) (bool, error) {
	// Get the user
	user, err := utils.InitializeRequest(req, *api.Client)
	if err != nil {
		// Bad user - pass the error through
		return false, err
	}

	// Run data validation
	if validation != nil {
		fmt.Println("Running field validation")
		err := validateFields(validation, item, nil, false)
		if err != nil {
			fmt.Println("Field validation error", err)
			return false, err
		}
	}

	// Marshal item into a dynamo map
	data, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Failed to marshal data", err)
		return false, err
	}

	// Build input
	input := dynamodb.PutItemInput{
		TableName: aws.String(api.DataTable),
		Item:      data,
	}

	// Get key schema
	output, err := api.Client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(api.DataTable),
	})
	if err != nil {
		fmt.Println("Failed to describe table", err)
		return false, nil
	}

	// Build filters
	conditions, hasConditions := filterbuilder.Filter(user, nil)

	// Append key schema conditions
	if !hasConditions {
		for _, schema := range output.Table.KeySchema {
			condition := expression.Name(*schema.AttributeName).AttributeNotExists()
			if !hasConditions {
				conditions = condition
				hasConditions = true
			} else {
				conditions = conditions.And(condition)
			}
		}
	}

	// Build expression
	expr, err := expression.NewBuilder().WithCondition(conditions).Build()
	if err != nil {
		fmt.Println("Encountered error while building expression", err)
		return false, err
	}

	// Update scan input
	input.ConditionExpression = expr.Condition()
	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()

	// Put the item into dynamo
	_, err = api.Client.PutItem(&input)
	if err != nil {
		fmt.Println("Error while attempting to add item to dynamo", err)
		return false, err
	}

	// Create audit log
	utils.AuditLog()

	return true, nil
}

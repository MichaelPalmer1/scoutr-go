package endpoints

import (
	"fmt"

	"github.com/MichaelPalmer1/simple-api-go/filterbuilder"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Update : Update an item
func (api *SimpleAPI) Update(req models.Request, partitionKey map[string]string, item map[string]string, validation map[string]FieldValidation) (bool, error) {
	// Get the user
	user, err := utils.InitializeRequest(req, *api.Client)
	if err != nil {
		// Bad user - pass the error through
		return false, err
	}

	// Run data validation
	if validation != nil {
		fmt.Println("Running field validation")
		err := validateFields(validation, item, nil, true)
		if err != nil {
			fmt.Println("Field validation error", err)
			return false, err
		}
	}

	// Build update expression
	var updateConds expression.UpdateBuilder
	hasUpdateConds := false
	for key, value := range item {
		if !hasUpdateConds {
			updateConds = expression.Set(expression.Name(key), expression.Value(value))
			hasUpdateConds = true
		} else {
			updateConds = updateConds.Set(expression.Name(key), expression.Value(value))
		}
	}

	// Build input
	input := dynamodb.UpdateItemInput{
		TableName: aws.String(api.DataTable),
	}

	// Build partition key
	dynamoKeyParts, err := dynamodbattribute.MarshalMap(partitionKey)
	if err != nil {
		fmt.Println("Failed to marshal partition key", err)
		return false, err
	}
	input.Key = dynamoKeyParts

	// Build filters
	var expr expression.Expression
	conditions, hasConditions := filterbuilder.Filter(user, nil)

	// Append key schema conditions
	if hasConditions {
		// Build expression
		expr, err = expression.NewBuilder().WithCondition(conditions).WithUpdate(updateConds).Build()
		if err != nil {
			fmt.Println("Encountered error while building expression", err)
			return false, err
		}

		// Update input
		input.ConditionExpression = expr.Condition()
	} else {
		expr, err = expression.NewBuilder().WithUpdate(updateConds).Build()
		if err != nil {
			fmt.Println("Encountered error while building expression", err)
			return false, err
		}
	}

	// Update input
	input.UpdateExpression = expr.Update()
	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()

	// Update the item in dynamo
	_, err = api.Client.UpdateItem(&input)
	if err != nil {
		fmt.Println("Error while attempting to update item in dynamo", err)
		return false, err
	}

	// Create audit log
	utils.AuditLog()

	return true, nil
}

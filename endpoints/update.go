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
func (api *SimpleAPI) Update(req models.Request, partitionKey map[string]string, item map[string]string, validation map[string]utils.FieldValidation) (interface{}, error) {
	var output interface{}

	// Get the user
	user, err := utils.InitializeRequest(req, *api.Client)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Run data validation
	if validation != nil {
		fmt.Println("Running field validation")
		err := utils.ValidateFields(validation, item, nil, true)
		if err != nil {
			fmt.Println("Field validation error", err)
			return nil, err
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
		TableName:    aws.String(api.DataTable),
		ReturnValues: aws.String("ALL_NEW"),
	}

	// Build partition key
	dynamoKeyParts, err := dynamodbattribute.MarshalMap(partitionKey)
	if err != nil {
		fmt.Println("Failed to marshal partition key", err)
		return nil, err
	}
	input.Key = dynamoKeyParts

	// Build filters
	var expr expression.Expression
	conditions, hasConditions := filterbuilder.Filter(user, nil)

	// Get key schema
	keySchema, err := api.Client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(api.DataTable),
	})
	if err != nil {
		fmt.Println("Failed to describe table", err)
		return nil, err
	}

	// Append key schema conditions
	for _, schema := range keySchema.Table.KeySchema {
		condition := expression.Name(*schema.AttributeName).AttributeExists()
		if !hasConditions {
			conditions = condition
			hasConditions = true
		} else {
			conditions = conditions.And(condition)
		}
	}

	// Build expression
	if hasConditions {
		expr, err = expression.NewBuilder().WithCondition(conditions).WithUpdate(updateConds).Build()
		if err != nil {
			fmt.Println("Encountered error while building expression", err)
			return nil, err
		}

		// Update input
		input.ConditionExpression = expr.Condition()
	} else {
		expr, err = expression.NewBuilder().WithUpdate(updateConds).Build()
		if err != nil {
			fmt.Println("Encountered error while building expression", err)
			return nil, err
		}
	}

	// Update input
	input.UpdateExpression = expr.Update()
	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()

	// Update the item in dynamo
	updatedItem, err := api.Client.UpdateItem(&input)
	if err != nil {
		fmt.Println("Error while attempting to update item in dynamo", err)

		// Check if this was a conditional check failure
		if _, ok := err.(*dynamodb.ConditionalCheckFailedException); ok {
			return nil, &models.BadRequest{
				Message: "Item does not exist or you do not have permission to update it",
			}
		}

		return nil, err
	}

	// Unmarshal into output interface
	dynamodbattribute.UnmarshalMap(updatedItem.Attributes, &output)

	// Create audit log
	utils.AuditLog()

	return output, nil
}

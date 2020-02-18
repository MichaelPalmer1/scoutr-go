package aws

import (
	"github.com/MichaelPalmer1/simple-api-go/lib/filtering"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/sirupsen/logrus"
)

// Update : Update an item
func (api *DynamoAPI) Update(req models.Request, partitionKey map[string]string, item map[string]string, validation map[string]utils.FieldValidation, auditAction string) (interface{}, error) {
	var output interface{}

	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Run data validation
	if validation != nil {
		log.Infoln("Running field validation")
		err := utils.ValidateFields(validation, item, nil, true)
		if err != nil {
			log.Errorln("Field validation error", err)
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
		TableName:    aws.String(api.Config.DataTable),
		ReturnValues: aws.String("ALL_NEW"),
	}

	// Build partition key
	dynamoKeyParts, err := dynamodbattribute.MarshalMap(partitionKey)
	if err != nil {
		log.Errorln("Failed to marshal partition key", err)
		return nil, err
	}
	input.Key = dynamoKeyParts

	// Build filters
	var expr expression.Expression
	conditions, hasConditions, err := filtering.Filter(user, nil)
	if err != nil {
		log.Errorln("Error encountered during filtering", err)
		return nil, err
	}

	// Get key schema
	keySchema, err := api.client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(api.Config.DataTable),
	})
	if err != nil {
		log.Errorln("Failed to describe table", err)
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
			log.Errorln("Encountered error while building expression", err)
			return nil, err
		}

		// Update input
		input.ConditionExpression = expr.Condition()
	} else {
		expr, err = expression.NewBuilder().WithUpdate(updateConds).Build()
		if err != nil {
			log.Errorln("Encountered error while building expression", err)
			return nil, err
		}
	}

	// Update input
	input.UpdateExpression = expr.Update()
	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()

	// Update the item in dynamo
	updatedItem, err := api.client.UpdateItem(&input)
	if err != nil {
		log.Errorln("Error while attempting to update item in dynamo", err)

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
	api.auditLog("UPDATE", req, *user, &partitionKey, &item)

	return output, nil
}

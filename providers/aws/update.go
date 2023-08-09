package aws

import (
	"context"
	"errors"

	"github.com/MichaelPalmer1/scoutr-go/models"
	"github.com/MichaelPalmer1/scoutr-go/providers/base"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	log "github.com/sirupsen/logrus"
)

// Update : Update an item
func (api DynamoAPI) Update(req models.Request, partitionKey map[string]string, item map[string]interface{}, validation map[string]models.FieldValidation, requiredFields []string, auditAction string) (interface{}, error) {
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
		err := api.ValidateFields(validation, requiredFields, item, nil)
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
		ReturnValues: types.ReturnValueAllNew,
	}

	// Build partition key
	dynamoKeyParts, err := attributevalue.MarshalMap(partitionKey)
	if err != nil {
		log.Errorln("Failed to marshal partition key", err)
		return nil, err
	}
	input.Key = dynamoKeyParts

	// Build filters
	var expr expression.Expression
	conditions, err := api.Filtering.Filter(user, nil, base.FilterActionUpdate)
	if err != nil {
		log.Errorln("Error encountered during filtering", err)
		return nil, err
	}

	// Get key schema
	keySchema, err := api.Client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(api.Config.DataTable),
	})
	if err != nil {
		log.Errorln("Failed to describe table", err)
		return nil, err
	}

	// Append key schema conditions
	for _, schema := range keySchema.Table.KeySchema {
		condition := expression.Name(*schema.AttributeName).AttributeExists()
		conditions = api.Filtering.And(conditions, condition)
	}

	// Build expression
	if conditions != nil {
		expr, err = expression.NewBuilder().WithCondition(conditions.(expression.ConditionBuilder)).WithUpdate(updateConds).Build()
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
	updatedItem, err := api.Client.UpdateItem(context.TODO(), &input)
	if err != nil {
		log.Errorln("Error while attempting to update item in dynamo", err)

		// Check if this was a conditional check failure
		var apiError smithy.APIError
		if errors.As(err, &apiError) && apiError.ErrorCode() == "ConditionalCheckFailedException" {
			return nil, &models.BadRequest{
				Message: "Item does not exist or you do not have permission to update it",
			}
		}

		return nil, err
	}

	// Unmarshal into output interface
	err = attributevalue.UnmarshalMap(updatedItem.Attributes, &output)
	if err != nil {
		return nil, err
	}

	// Create audit log
	api.auditLog(auditAction, req, *user, &partitionKey, &item)

	return output, nil
}

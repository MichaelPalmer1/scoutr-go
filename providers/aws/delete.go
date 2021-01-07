package aws

import (
	"github.com/MichaelPalmer1/scoutr-go/models"
	"github.com/MichaelPalmer1/scoutr-go/providers/base"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/sirupsen/logrus"
)

// Delete : Delete an item
func (api DynamoAPI) Delete(req models.Request, partitionKey map[string]string) error {
	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return err
	}

	// Build input
	input := dynamodb.DeleteItemInput{
		TableName:    aws.String(api.Config.DataTable),
		ReturnValues: aws.String("ALL_OLD"),
	}

	// Build partition key
	dynamoKeyParts, err := dynamodbattribute.MarshalMap(partitionKey)
	if err != nil {
		log.Errorln("Failed to marshal partition key", err)
		return err
	}
	input.Key = dynamoKeyParts

	// Build filters
	var expr expression.Expression
	rawConds, err := api.Filtering.Filter(user, nil, base.FilterActionDelete)
	if err != nil {
		log.Errorln("Error encountered during filtering", err)
		return err
	}

	// Cast to condition builder
	var conditions interface{}
	if rawConds != nil {
		conditions = rawConds.(expression.ConditionBuilder)
	}

	// Get key schema
	keySchema, err := api.Client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(api.Config.DataTable),
	})
	if err != nil {
		log.Errorln("Failed to describe table", err)
		return err
	}

	// Append key schema conditions
	for _, schema := range keySchema.Table.KeySchema {
		var condition interface{}
		condition = expression.Name(*schema.AttributeName).AttributeExists()
		conditions = api.Filtering.And(conditions, condition)
	}

	// Build expression
	if rawConds != nil {
		expr, err = expression.NewBuilder().WithCondition(conditions).Build()
		if err != nil {
			log.Errorln("Encountered error while building expression", err)
			return err
		}

		// Update input
		input.ConditionExpression = expr.Condition()
		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
	}

	// Delete the item from dynamo
	_, err = api.Client.DeleteItem(&input)
	if err != nil {
		log.Errorln("Error while attempting to delete item in dynamo", err)

		// Check if this was a conditional check failure
		if _, ok := err.(*dynamodb.ConditionalCheckFailedException); ok {
			return &models.BadRequest{
				Message: "Item does not exist or you do not have permission to delete it",
			}
		}

		return err
	}

	// Create audit log
	api.auditLog("DELETE", req, *user, &partitionKey, nil)

	return nil
}

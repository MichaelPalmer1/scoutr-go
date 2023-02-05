package aws

import (
	"context"
	"errors"

	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go"
	log "github.com/sirupsen/logrus"
)

// Delete : Delete an item
func (api DynamoAPI) Delete(request types.Request, partitionKey map[string]interface{}) error {
	// Get the user
	user, err := api.InitializeRequest(request)
	if err != nil {
		// Bad user - pass the error through
		return err
	}

	// Build partition key
	dynamoKeyParts, err := attributevalue.MarshalMap(partitionKey)
	if err != nil {
		log.Errorln("Failed to marshal partition key", err)
		return err
	}

	// Build filters
	var expr expression.Expression
	conditions, err := api.filtering.Filter(user, nil, base.FilterActionDelete)
	if err != nil {
		log.Errorln("Error encountered during filtering", err)
		return err
	}

	// Get key schema
	keySchema, err := api.Client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(api.Config.DataTable),
	})
	if err != nil {
		log.Errorln("Failed to describe table", err)
		return err
	}

	// Append key schema conditions
	for _, schema := range keySchema.Table.KeySchema {
		conditions = api.filtering.And(conditions, expression.Name(*schema.AttributeName).AttributeExists())
	}

	// Build expression
	if conditions != nil {
		expr, err = expression.NewBuilder().WithCondition(conditions.(expression.ConditionBuilder)).Build()
		if err != nil {
			log.Errorln("Encountered error while building expression", err)
			return err
		}
	}

	// Delete the item from dynamo
	if err := api.DeleteItem(api.Config.DataTable, dynamoKeyParts, &expr); err != nil {
		log.Errorln("Error while attempting to delete item in dynamo", err)

		// Check if this was a conditional check failure
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ConditionalCheckFailedException" {
			return &types.BadRequest{
				Message: "Item does not exist or you do not have permission to delete it",
			}
		}

		return err
	}

	// Create audit log
	api.auditLog(base.AuditActionDelete, request, user, partitionKey, nil)

	return nil
}

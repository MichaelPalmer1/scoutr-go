package aws

import (
	"context"
	"errors"

	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	log "github.com/sirupsen/logrus"
)

// Create : Create an item
func (api DynamoAPI) Create(req types.Request, item map[string]interface{}, validation map[string]types.FieldValidation, requiredFields []string) error {
	var conditions interface{}

	// Get the user
	user, err := api.PrepareCreate(req, item, validation, requiredFields)
	if err != nil {
		// Bad user - pass the error through
		return err
	}

	// Get key schema
	output, err := api.Client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(api.Config.DataTable),
	})
	if err != nil {
		log.Errorln("Failed to describe table", err)
		return err
	}

	// Append key schema conditions
	partitionKey := ""
	for _, schema := range output.Table.KeySchema {
		if schema.KeyType == dynamoTypes.KeyTypeHash {
			partitionKey = *schema.AttributeName
		}
		conditions = api.filtering.And(conditions, expression.Name(*schema.AttributeName).AttributeNotExists())
	}

	// Build expression
	expr, err := expression.NewBuilder().WithCondition(conditions.(expression.ConditionBuilder)).Build()
	if err != nil {
		log.Errorln("Encountered error while building expression", err)
		return err
	}

	// Put the item into the table
	if err := api.PutItem(api.Config.DataTable, item, &expr); err != nil {
		log.WithError(err).Errorln("Encountered error while attempting to create record")

		// Check if this was a conditional check failure
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ConditionalCheckFailedException" {
			return &types.BadRequest{
				Message: "Item already exists or you do not have permission to create it",
			}
		}

		return err
	}

	// Create audit log
	api.auditLog(base.AuditActionCreate, req, user, map[string]interface{}{partitionKey: item[partitionKey]}, nil)

	return nil
}

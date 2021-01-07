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

// Create : Create an item
func (api DynamoAPI) Create(req models.Request, item map[string]interface{}, validation map[string]models.FieldValidation, requiredFields []string) error {
	var conditions interface{}

	// Get the user
	user, err := api.PrepareCreate(req, item, validation, requiredFields)
	if err != nil {
		// Bad user - pass the error through
		return err
	}

	// Get key schema
	output, err := api.Client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(api.Config.DataTable),
	})
	if err != nil {
		log.Errorln("Failed to describe table", err)
		return err
	}

	// Append key schema conditions
	partitionKey := ""
	for _, schema := range output.Table.KeySchema {
		if *schema.KeyType == "HASH" {
			partitionKey = *schema.AttributeName
		}
		conditions = api.Filtering.And(conditions, expression.Name(*schema.AttributeName).AttributeNotExists())
	}

	// Marshal item into a dynamo map
	data, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Errorln("Failed to marshal data", err)
		return err
	}

	// Build expression
	expr, err := expression.NewBuilder().WithCondition(conditions.(expression.ConditionBuilder)).Build()
	if err != nil {
		log.Errorln("Encountered error while building expression", err)
		return err
	}

	// Build input
	input := dynamodb.PutItemInput{
		TableName:                 aws.String(api.Config.DataTable),
		Item:                      data,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Put the item into the table
	_, err = api.Client.PutItem(&input)
	if err != nil {
		log.WithError(err).Errorln("Encountered error while attempting to create record")

		// Check if this was a conditional check failure
		if _, ok := err.(*dynamodb.ConditionalCheckFailedException); ok {
			return &models.BadRequest{
				Message: "Item already exists or you do not have permission to create it",
			}
		}

		return err
	}

	// Create audit log
	api.auditLog(base.AuditActionCreate, req, *user, &map[string]string{partitionKey: item[partitionKey].(string)}, nil)

	return nil
}

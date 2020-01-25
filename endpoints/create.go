package endpoints

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

// Create : Create an item
func (api *SimpleAPI) Create(req models.Request, item map[string]string, validation map[string]utils.FieldValidation) error {
	// Get the user
	user, err := api.initializeRequest(req, *api.Client)
	if err != nil {
		// Bad user - pass the error through
		return err
	}

	// Run data validation
	if validation != nil {
		log.Infoln("Running field validation")
		err := utils.ValidateFields(validation, item, nil, false)
		if err != nil {
			log.Errorln("Field validation error", err)
			return err
		}
	}

	// Marshal item into a dynamo map
	data, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Errorln("Failed to marshal data", err)
		return err
	}

	// Build input
	input := dynamodb.PutItemInput{
		TableName: aws.String(api.Config.DataTable),
		Item:      data,
	}

	// Get key schema
	output, err := api.Client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(api.Config.DataTable),
	})
	if err != nil {
		log.Errorln("Failed to describe table", err)
		return err
	}

	// Build filters
	conditions, hasConditions := filtering.Filter(user, nil)

	// Append key schema conditions
	partitionKey := ""
	for _, schema := range output.Table.KeySchema {
		if *schema.KeyType == "HASH" {
			partitionKey = *schema.AttributeName
		}
		condition := expression.Name(*schema.AttributeName).AttributeNotExists()
		if !hasConditions {
			conditions = condition
			hasConditions = true
		} else {
			conditions = conditions.And(condition)
		}
	}

	// Build expression
	expr, err := expression.NewBuilder().WithCondition(conditions).Build()
	if err != nil {
		log.Errorln("Encountered error while building expression", err)
		return err
	}

	// Update input
	input.ConditionExpression = expr.Condition()
	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()

	// Put the item into dynamo
	_, err = api.Client.PutItem(&input)
	if err != nil {
		log.Errorln("Error while attempting to add item to dynamo", err)

		// Check if this was a conditional check failure
		if _, ok := err.(*dynamodb.ConditionalCheckFailedException); ok {
			return &models.BadRequest{
				Message: "Item already exists or you do not have permission to create it",
			}
		}

		return err
	}

	// Create audit log
	api.auditLog("CREATE", req, *user, &map[string]string{partitionKey: item[partitionKey]}, nil)

	return nil
}

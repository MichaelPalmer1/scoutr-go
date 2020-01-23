package endpoints

import (
	"fmt"

	"github.com/MichaelPalmer1/simple-api-go/filterbuilder"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Get : Get an item from the table
func (api *SimpleAPI) Get(req models.Request, id string) (models.Record, error) {
	var partitionKey string

	// Get the user
	user, err := api.initializeRequest(req, *api.Client)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Lookup the partition key
	tableInfo, err := api.Client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(api.Config.DataTable),
	})
	if err != nil {
		fmt.Println("Failed to describe table", err)
		return nil, err
	}

	// Get partition key
	for _, schema := range tableInfo.Table.KeySchema {
		if *schema.KeyType == "HASH" {
			partitionKey = *schema.AttributeName
			break
		}
	}

	// Build filters
	conditions, hasConditions := filterbuilder.Filter(user, map[string]string{})
	keyCondition := expression.Name(partitionKey).Equal(expression.Value(id))
	if hasConditions {
		conditions = conditions.And(keyCondition)
	} else {
		conditions = keyCondition
	}

	// Build expression
	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		return nil, err
	}

	// Build scan input
	input := dynamodb.ScanInput{
		TableName:                 aws.String(api.Config.DataTable),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Download the data
	data, err := scan(&input, api.Client)
	if err != nil {
		fmt.Println("Error while attempting to list records", err)
		return nil, nil
	}

	// Filter the response
	utils.PostProcess(data, user)

	// Make sure only a single record was returned
	if len(data) > 1 {
		return nil, &models.BadRequest{
			Message: "Multiple items returned",
		}
	} else if len(data) == 0 {
		return nil, &models.NotFound{
			Message: "Item does not exist or you do not have permission to view it",
		}
	}

	// TODO: Create audit log
	utils.AuditLog()

	return data[0], nil
}

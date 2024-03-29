package aws

import (
	"context"

	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
)

// Get : Get an item from the table
func (api DynamoAPI) Get(req types.Request, id string) (types.Record, error) {
	var partitionKey string

	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	// Lookup the partition key
	tableInfo, err := api.Client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(api.Config.DataTable),
	})
	if err != nil {
		log.Errorln("Failed to describe table", err)
		return nil, err
	}

	// Get partition key
	for _, schema := range tableInfo.Table.KeySchema {
		if schema.KeyType == dynamoTypes.KeyTypeHash {
			partitionKey = *schema.AttributeName
			break
		}
	}

	// Build filters
	conditions, err := api.filtering.Filter(user, nil, "")
	if err != nil {
		log.Errorln("Error encountered during filtering", err)
		return nil, err
	}

	// Build key condition
	keyCondition := expression.Name(partitionKey).Equal(expression.Value(id))
	conditions = api.filtering.And(conditions, keyCondition)

	// Build expression
	expr, err := expression.NewBuilder().WithFilter(conditions.(expression.ConditionBuilder)).Build()
	if err != nil {
		return nil, err
	}

	// Build scan input
	input := &dynamodb.ScanInput{
		TableName:                 aws.String(api.Config.DataTable),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Download the data
	data, err := Scan[types.Record](api.Client, input)
	if err != nil {
		log.Errorln("Error while attempting to list records", err)
		return nil, err
	}

	// Filter the response
	api.PostProcess(data, user)

	// Make sure only a single record was returned
	if len(data) > 1 {
		return nil, &types.BadRequest{
			Message: "Multiple items returned",
		}
	} else if len(data) == 0 {
		return nil, &types.NotFound{
			Message: "Item does not exist or you do not have permission to view it",
		}
	}

	// Create audit log
	api.auditLog(base.AuditActionGet, req, user, map[string]interface{}{partitionKey: id}, nil)

	return data[0], nil
}

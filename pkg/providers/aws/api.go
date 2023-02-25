package aws

import (
	"context"

	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtraildata"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cenkalti/backoff/v4"
)

// DynamoAPI : API, based off of Scoutr, used to talk to AWS DynamoDB
type DynamoAPI struct {
	*base.Scoutr
	Client           types.DynamoClientAPI
	filtering        DynamoFiltering
	auditClient      *cloudtraildata.Client
	cloudTrailClient *cloudtrail.Client
}

// Init : Initialize the Dynamo client
func (api *DynamoAPI) Init(config aws.Config) {
	api.Client = dynamodb.NewFromConfig(config)
	api.auditClient = cloudtraildata.NewFromConfig(config)
	api.cloudTrailClient = cloudtrail.NewFromConfig(config)
	api.filtering = NewFilter()
	api.ScoutrBase = api
}

func Scan[T any](client types.DynamoClientAPI, input *dynamodb.ScanInput) ([]T, error) {
	var results []T
	paginator := dynamodb.NewScanPaginator(client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		var data []T
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &data); err != nil {
			return nil, err
		} else {
			results = append(results, data...)
		}
	}

	return results, nil
}

func Query[T any](client types.DynamoClientAPI, input *dynamodb.QueryInput) ([]T, error) {
	var results []T
	paginator := dynamodb.NewQueryPaginator(client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		var data []T
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &data); err != nil {
			return nil, err
		} else {
			results = append(results, data...)
		}
	}

	return results, nil
}

func GetItem[T any](client types.DynamoClientAPI, input *dynamodb.GetItemInput) (*T, error) {
	var output *T
	var item map[string]dynamoTypes.AttributeValue

	// Backoff operation
	fn := func() error {
		result, err := client.GetItem(context.TODO(), input)
		if err != nil {
			return err
		}

		item = result.Item

		return nil
	}

	// Perform exponential backoff
	if err := backoff.Retry(fn, backoff.NewExponentialBackOff()); err != nil {
		return nil, err
	}

	// Item does not exist
	if item == nil {
		return nil, nil
	}

	if err := attributevalue.UnmarshalMap(item, &output); err != nil {
		return nil, err
	}

	return output, nil
}

func (api *DynamoAPI) PutItem(table string, item interface{}, expr *expression.Expression) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item:      av,
	}

	if expr != nil {
		input.ConditionExpression = expr.Condition()
		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
	}

	if _, err := api.Client.PutItem(context.TODO(), input); err != nil {
		return err
	}

	return nil
}

func UpdateItem[T any](client types.DynamoClientAPI, table string, key map[string]dynamoTypes.AttributeValue, expr expression.Expression) (*T, error) {
	var output *T

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(table),
		Key:                       key,
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
		ReturnValues:              dynamoTypes.ReturnValueAllNew,
	}

	if result, err := client.UpdateItem(context.TODO(), input); err != nil {
		return nil, err
	} else if err := attributevalue.UnmarshalMap(result.Attributes, &output); err != nil {
		return nil, err
	}

	return output, nil
}

func (api *DynamoAPI) DeleteItem(table string, key map[string]dynamoTypes.AttributeValue, expr *expression.Expression) error {
	input := &dynamodb.DeleteItemInput{
		TableName:    aws.String(table),
		Key:          key,
		ReturnValues: dynamoTypes.ReturnValueAllOld,
	}

	if expr != nil {
		input.ConditionExpression = expr.Condition()
		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
	}

	if _, err := api.Client.DeleteItem(context.TODO(), input); err != nil {
		return err
	}

	return nil
}

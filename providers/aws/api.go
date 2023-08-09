package aws

import (
	"context"

	"github.com/MichaelPalmer1/scoutr-go/models"
	"github.com/MichaelPalmer1/scoutr-go/providers/base"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoAPI : API, based off of Scoutr, used to talk to AWS DynamoDB
type DynamoAPI struct {
	*base.Scoutr
	Client    *dynamodb.Client
	Filtering DynamoFiltering
}

// Init : Initialize the Dynamo client
func (api *DynamoAPI) Init(config aws.Config) {
	api.Client = dynamodb.NewFromConfig(config)
	f := DynamoFiltering{}
	f.Filtering = &base.Filtering{
		FilterBase:    &f,
		ScoutrFilters: &f,
	}
	api.Filtering = f
	api.ScoutrBase = api
}

func scan(input *dynamodb.ScanInput, client *dynamodb.Client) ([]models.Record, error) {
	results := []models.Record{}
	paginator := dynamodb.NewScanPaginator(client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		// Unmarshal data into Record model
		records := []models.Record{}
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &records); err != nil {
			return nil, err
		}

		// Append records to results
		results = append(results, records...)
	}

	return results, nil
}

func scanAudit(input *dynamodb.ScanInput, client *dynamodb.Client) ([]models.AuditLog, error) {
	results := []models.AuditLog{}
	paginator := dynamodb.NewScanPaginator(client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		// Unmarshal data into AuditLog model
		records := []models.AuditLog{}
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &records); err != nil {
			return nil, err
		}

		// Append records to results
		results = append(results, records...)
	}

	return results, nil
}

func (api *DynamoAPI) storeItem(table string, item map[string]interface{}) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item:      av,
	}

	_, err = api.Client.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

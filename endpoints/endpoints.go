package endpoints

import (
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// SimpleAPI : Class
type SimpleAPI struct {
	DataTable string
	Client    *dynamodb.DynamoDB
}

func scan(input *dynamodb.ScanInput, client *dynamodb.DynamoDB) ([]models.Record, error) {
	results, err := client.Scan(input)
	if err != nil {
		return nil, err
	}

	records := []models.Record{}
	dynamodbattribute.UnmarshalListOfMaps(results.Items, &records)

	return records, nil
}

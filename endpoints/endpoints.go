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
	results := []models.Record{}
	err := client.ScanPages(input,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			// Unmarshal data into Record model
			records := []models.Record{}
			dynamodbattribute.UnmarshalListOfMaps(page.Items, &records)

			// Append to output
			results = append(results, records...)

			return true
		},
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}

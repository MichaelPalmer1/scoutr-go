package aws

import (
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/providers/base"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoAPI : API, based off of SimpleAPI, used to talk to AWS DynamoDB
type DynamoAPI struct {
	*base.SimpleAPI
	Client    *dynamodb.DynamoDB
	Filtering DynamoFiltering
}

// Init : Initialize the Dynamo client
func (api *DynamoAPI) Init(config *aws.Config) {
	sess := session.Must(session.NewSession(config))
	api.Client = dynamodb.New(sess)
}

func scan(input *dynamodb.ScanInput, client *dynamodb.DynamoDB) ([]models.Record, error) {
	results := []models.Record{}
	var lastErr error
	err := client.ScanPages(input,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			// Unmarshal data into Record model
			records := []models.Record{}
			err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &records)
			if err != nil {
				lastErr = err
				return false
			}

			// Append records to results
			results = append(results, records...)

			return true
		},
	)
	if err != nil {
		return nil, err
	}
	if lastErr != nil {
		return nil, lastErr
	}

	return results, nil
}

func scanAudit(input *dynamodb.ScanInput, client *dynamodb.DynamoDB) ([]models.AuditLog, error) {
	results := []models.AuditLog{}
	var lastErr error
	err := client.ScanPages(input,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			// Unmarshal data into AuditLog model
			records := []models.AuditLog{}
			err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &records)
			if err != nil {
				lastErr = err
				return false
			}

			// Append records to results
			results = append(results, records...)

			return true
		},
	)
	if err != nil {
		return nil, err
	}
	if lastErr != nil {
		return nil, lastErr
	}

	return results, nil
}

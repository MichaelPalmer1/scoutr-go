package aws

import (
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/providers/base"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/sirupsen/logrus"
)

type DynamoAPI struct {
	base.SimpleAPI
	client *dynamodb.DynamoDB
}

// Init : Initialize the Dynamo client
func (api *DynamoAPI) Init(config aws.Config) {
	sess := session.Must(session.NewSession(&config))
	api.client = dynamodb.New(sess)
}

func (api *DynamoAPI) List(req models.Request) ([]models.Record, error) {
	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	input := dynamodb.ScanInput{
		TableName: aws.String(api.Config.DataTable),
	}

	// Copy queryParams into params
	params := make(map[string]string)
	for key, value := range req.QueryParams {
		params[key] = value
	}

	// Merge pathParams into params
	for key, value := range req.PathParams {
		params[key] = value
	}

	// Generate dynamic search
	searchKey, hasSearchKey := req.PathParams["search_key"]
	searchValue, hasSearchValue := req.PathParams["search_value"]
	if hasSearchKey && hasSearchValue {
		// Map the search key and value into path params
		params[searchKey] = searchValue
		delete(params, "search_key")
		delete(params, "search_value")
	}

	// Build filters
	conditions, hasConditions, err := Filter(user, params)
	if err != nil {
		log.Errorln("Error encountered during filtering", err)
		return nil, err
	}
	if hasConditions {
		expr, err := expression.NewBuilder().WithFilter(conditions).Build()
		if err != nil {
			return nil, err
		}

		// Update scan input
		input.FilterExpression = expr.Filter()
		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
	}

	// Download the data
	data, err := scan(&input, api.client)
	if err != nil {
		log.Errorln("Error while attempting to list records", err)
		return nil, nil
	}

	// Filter the response
	utils.PostProcess(data, user)

	// Create audit log
	// api.auditLog("LIST", req, *user, nil, nil)

	return data, nil
}

func scan(input *dynamodb.ScanInput, client *dynamodb.DynamoDB) ([]models.Record, error) {
	results := []models.Record{}
	err := client.ScanPages(input,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			// Unmarshal data into Record model
			records := []models.Record{}
			dynamodbattribute.UnmarshalListOfMaps(page.Items, &records)

			// Append records to results
			results = append(results, records...)

			return true
		},
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}

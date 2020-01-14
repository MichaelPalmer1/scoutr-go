package endpoints

import (
	"fmt"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/filterbuilder"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/aws"
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

// ListTable : Lists all items in a table
func (api *SimpleAPI) ListTable(req models.Request, uniqueKey string, pathParams map[string]string, queryParams map[string]string) []models.Record {
	// Get the user
	user := utils.InitializeRequest(req, *api.Client)
	if user == nil {
		fmt.Println("BAD USER")
		return nil
	}

	input := dynamodb.ScanInput{
		TableName: aws.String(api.DataTable),
	}

	// Generate dynamic search
	searchKey, hasSearchKey := pathParams["search_key"]
	searchValue, hasSearchValue := pathParams["search_value"]
	if hasSearchKey && hasSearchValue {
		// Map the search key and value into path params
		pathParams[searchKey] = searchValue
		delete(pathParams, "search_key")
		delete(pathParams, "search_value")
	}

	// Merge pathParams into queryParams
	for key, value := range pathParams {
		queryParams[key] = value
	}

	// Build filters
	conditions, err := filterbuilder.Filter(user, queryParams)

	// Update scan input
	input.FilterExpression = conditions.Filter()
	input.ProjectionExpression = conditions.Projection()
	input.ExpressionAttributeNames = conditions.Names()
	input.ExpressionAttributeValues = conditions.Values()

	// Download the data
	data, err := scan(&input, api.Client)
	if err != nil {
		fmt.Println("Error while attempting to list records", err)
		return nil
	}

	// Filter the response
	fmt.Println(user)
	fmt.Println(data)
	filteredData := utils.PostProcess(data, user)
	fmt.Println(filteredData)

	// Sort the response if unique key was specified


	// Create audit log
	utils.AuditLog()

	return data
}

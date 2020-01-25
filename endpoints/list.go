package endpoints

import (
	"github.com/MichaelPalmer1/simple-api-go/lib/filtering"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/sirupsen/logrus"
)

// ListTable : Lists all items in a table
func (api *SimpleAPI) ListTable(req models.Request, uniqueKey string, pathParams map[string]string, queryParams map[string]string) ([]models.Record, error) {
	// Get the user
	user, err := api.initializeRequest(req, *api.Client)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	input := dynamodb.ScanInput{
		TableName: aws.String(api.Config.DataTable),
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
	conditions, hasConditions := filtering.Filter(user, queryParams)
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
	data, err := scan(&input, api.Client)
	if err != nil {
		log.Errorln("Error while attempting to list records", err)
		return nil, nil
	}

	// Filter the response
	utils.PostProcess(data, user)

	// TODO: Sort the response if unique key was specified

	// Create audit log
	api.auditLog("LIST", req, *user, nil, nil)

	return data, nil
}

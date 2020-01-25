package simpleapi

import (
	"sort"

	"github.com/MichaelPalmer1/simple-api-go/lib/filtering"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/sirupsen/logrus"
)

// ListTable : Lists all items in a table
func (api *SimpleAPI) ListTable(req models.Request) ([]models.Record, error) {
	// Get the user
	user, err := api.initializeRequest(req, *api.Client)
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
	conditions, hasConditions, err := filtering.Filter(user, params)
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
	data, err := scan(&input, api.Client)
	if err != nil {
		log.Errorln("Error while attempting to list records", err)
		return nil, nil
	}

	// Filter the response
	utils.PostProcess(data, user)

	// Create audit log
	api.auditLog("LIST", req, *user, nil, nil)

	return data, nil
}

// ListUniqueValues : Lists unique values in a table
func (api *SimpleAPI) ListUniqueValues(req models.Request, uniqueKey string) ([]string, error) {
	// Get the user
	user, err := api.initializeRequest(req, *api.Client)
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

	// Build filters
	conditions, hasConditions, err := filtering.Filter(user, params)
	if err != nil {
		log.Errorln("Error encountered during filtering", err)
		return nil, err
	}

	// Build unique key condition
	condition := expression.Name(uniqueKey).AttributeExists()

	// Append unique key condition
	if hasConditions {
		conditions = conditions.And(condition)
	} else {
		conditions = condition
	}

	// Build projection expression
	projection := expression.NamesList(expression.Name(uniqueKey))

	// Build filter
	expr, err := expression.NewBuilder().WithFilter(conditions).WithProjection(projection).Build()
	if err != nil {
		return nil, err
	}

	// Update scan input
	input.FilterExpression = expr.Filter()
	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()
	input.ProjectionExpression = expr.Projection()

	// Download the data
	data, err := scan(&input, api.Client)
	if err != nil {
		log.Errorln("Error while attempting to list records", err)
		return nil, nil
	}

	// Filter the response
	utils.PostProcess(data, user)

	// Find unique values
	var values []string
	hasValue := make(map[string]bool)
	for _, item := range data {
		_, ok := hasValue[item[uniqueKey].(string)]
		if !ok {
			values = append(values, item[uniqueKey].(string))
			hasValue[item[uniqueKey].(string)] = true
		}
	}

	// Sort the data
	sort.Strings(values)

	// Create audit log
	api.auditLog("LIST", req, *user, nil, nil)

	return values, nil
}

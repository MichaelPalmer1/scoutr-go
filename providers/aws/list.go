package aws

import (
	"sort"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/sirupsen/logrus"
)

// List : Lists all items in a table
func (api DynamoAPI) List(req models.Request) ([]models.Record, error) {
	// Get the user
	user, err := api.InitializeRequest(api, req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	input := dynamodb.ScanInput{
		TableName: aws.String(api.Config.DataTable),
	}

	// Build filters
	rawConds, hasConditions, err := api.Filter(&api.Filtering, user, api.BuildParams(req))
	if err != nil {
		log.Errorln("Error encountered during filtering", err)
		return nil, err
	}
	if hasConditions {
		conditions := rawConds.(expression.ConditionBuilder)
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
		return nil, err
	}

	// Filter the response
	api.PostProcess(data, user)

	// Create audit log
	api.auditLog("LIST", req, *user, nil, nil)

	return data, nil
}

// ListUniqueValues : Lists unique values in a table
func (api DynamoAPI) ListUniqueValues(req models.Request, uniqueKey string) ([]string, error) {
	// Get the user
	user, err := api.InitializeRequest(api, req)
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
	rawConds, hasConditions, err := api.Filter(&api.Filtering, user, params)
	if err != nil {
		log.Errorln("Error encountered during filtering", err)
		return nil, err
	}

	// Cast to condition builder
	var conditions expression.ConditionBuilder
	if hasConditions {
		conditions = rawConds.(expression.ConditionBuilder)
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
		return nil, err
	}

	// Filter the response
	api.PostProcess(data, user)

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

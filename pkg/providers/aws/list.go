package aws

import (
	"sort"

	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sirupsen/logrus"
)

// List : Lists all items in a table
func (api DynamoAPI) List(req types.Request) ([]types.Record, error) {
	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(api.Config.DataTable),
	}

	// Build filters
	conditions, err := api.filtering.Filter(user, api.BuildParams(req), "")
	if err != nil {
		logrus.WithError(err).Error("Filtering failed")

		if api.Config.ErrorFunc != nil {
			api.Config.ErrorFunc(&req, user, err)
		}

		return nil, err
	}
	if conditions != nil {
		expr, err := expression.NewBuilder().WithFilter(conditions.(expression.ConditionBuilder)).Build()
		if err != nil {
			logrus.WithError(err).Error("Failed to build filter expression")

			if api.Config.ErrorFunc != nil {
				api.Config.ErrorFunc(&req, user, err)
			}

			return nil, err
		}

		// Update scan input
		input.FilterExpression = expr.Filter()
		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
	}

	// Download the data
	data, err := Scan[types.Record](api.Client, input)
	if err != nil {
		logrus.WithError(err).Error("Failed to list records")

		if api.Config.ErrorFunc != nil {
			api.Config.ErrorFunc(&req, user, err)
		}

		return nil, err
	}

	// Filter the response
	api.PostProcess(data, user)

	// Create audit log
	api.auditLog(base.AuditActionList, req, user, nil, nil)

	return data, nil
}

// ListUniqueValues : Lists unique values in a table
func (api DynamoAPI) ListUniqueValues(req types.Request, uniqueKey string) ([]string, error) {
	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(api.Config.DataTable),
	}

	// Copy queryParams into params
	params := make(map[string][]string)
	for key, values := range req.QueryParams {
		params[key] = append(params[key], values...)
	}

	// Build filters
	conditions, err := api.filtering.Filter(user, params, "")
	if err != nil {
		logrus.WithError(err).Error("Filtering failed")

		if api.Config.ErrorFunc != nil {
			api.Config.ErrorFunc(&req, user, err)
		}

		return nil, err
	}

	// Build unique key condition
	conditions = api.filtering.And(conditions, expression.Name(uniqueKey).AttributeExists())

	// Build projection expression
	projection := expression.NamesList(expression.Name(uniqueKey))

	// Build filter
	expr, err := expression.NewBuilder().WithFilter(conditions.(expression.ConditionBuilder)).WithProjection(projection).Build()
	if err != nil {
		logrus.WithError(err).Error("Failed to build filter expression")

		if api.Config.ErrorFunc != nil {
			api.Config.ErrorFunc(&req, user, err)
		}

		return nil, err
	}

	// Update scan input
	input.FilterExpression = expr.Filter()
	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()
	input.ProjectionExpression = expr.Projection()

	// Download the data
	data, err := Scan[types.Record](api.Client, input)
	if err != nil {
		logrus.WithError(err).Error("Failed to list records")

		if api.Config.ErrorFunc != nil {
			api.Config.ErrorFunc(&req, user, err)
		}

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
	api.auditLog(base.AuditActionList, req, user, nil, nil)

	return values, nil
}

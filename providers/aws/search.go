package aws

import (
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/sirupsen/logrus"
)

// Search : Search items in the table
func (api DynamoAPI) Search(req models.Request, key string, values []string) ([]models.Record, error) {
	// Get the user
	user, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	input := dynamodb.ScanInput{
		TableName: aws.String(api.Config.DataTable),
	}

	// Build filters
	rawConds, err := api.MultiFilter(&api.Filtering, user, key, values)
	if err != nil {
		log.Errorln("Error encountered during filtering", err)
		return nil, err
	}
	conditions := rawConds.(expression.ConditionBuilder)
	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		log.Errorln("Failed to build expression", err)
	}

	// Update scan input
	input.FilterExpression = expr.Filter()
	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()

	// Download the data
	data, err := scan(&input, api.Client)
	if err != nil {
		log.Errorln("Error while attempting to list records", err)
		return nil, err
	}

	// Filter the response
	api.PostProcess(data, user)

	// Create audit log
	api.auditLog("SEARCH", req, *user, nil, nil)

	return data, nil
}

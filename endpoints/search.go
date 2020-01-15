package endpoints

import (
	"fmt"

	"github.com/MichaelPalmer1/simple-api-go/filterbuilder"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Search : Search items in the table
func (api *SimpleAPI) Search(req models.Request, key string, values []string) ([]models.Record, error) {
	// Get the user
	user, err := utils.InitializeRequest(req, *api.Client)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	input := dynamodb.ScanInput{
		TableName: aws.String(api.DataTable),
	}

	// Build filters
	conditions := filterbuilder.MultiFilter(user, key, values)
	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		fmt.Println("Failed to build expression", err)
	}

	// Update scan input
	input.FilterExpression = expr.Filter()
	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()

	// Download the data
	data, err := scan(&input, api.Client)
	if err != nil {
		fmt.Println("Error while attempting to list records", err)
		return nil, nil
	}

	// Filter the response
	fmt.Println(user)
	fmt.Println(data)
	filteredData := utils.PostProcess(data, user)
	fmt.Println(filteredData)

	// Sort the response if unique key was specified

	// Create audit log
	utils.AuditLog()

	return data, nil
}

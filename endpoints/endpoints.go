package endpoints

import (
	"fmt"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// SimpleAPI : Class
type SimpleAPI struct {
	DataTable string
	Client    dynamodb.DynamoDB
}

func scan(tableName string, client *dynamodb.DynamoDB) (interface{}, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	results, err := client.Scan(input)
	if err != nil {
		return nil, err
	}

	var output interface{}
	dynamodbattribute.UnmarshalListOfMaps(results.Items, output)

	return output, nil
}

func buildFilters(user *models.User) *expression.Expression {
	return nil
}

// ListTable : Lists all items in a table
func (api *SimpleAPI) ListTable(req models.Request, uniqueKey string, pathParams map[string]string, queryParams map[string]string) interface{} {
	// Get the user
	// user := utils.(req)
	user := utils.InitializeRequest(req)

	// Build filters
	buildFilters(&user)

	// Download the data
	data, err := scan(api.DataTable, &api.Client)
	if err != nil {
		fmt.Println("Error while attempting to list records", err)
		return nil
	}

	// Filter the response

	// Sort the response if unique key was specified

	// Create audit log
	utils.AuditLog()

	return data
}

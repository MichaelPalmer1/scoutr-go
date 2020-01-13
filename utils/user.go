package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// GetUser : Get user data
func GetUser(id string, authTable string, client *dynamodb.DynamoDB, userData map[string]string, groups []string) (interface{}, error) {
	result, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(authTable),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	})
	if err != nil {
		fmt.Println("Failed to get user", err)
		return nil, err
	}

	fmt.Println(result)

	return result, nil
}

package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// GetUser : Get user data
func GetUser(id string, authTable string, groupTable string, client *dynamodb.DynamoDB, userData map[string]string, groups []string) (*models.User, error) {
	var user models.User

	// Get the user from Dynamo
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
	
	// Unmarshal result into User object
	dynamodbattribute.UnmarshalMap(result.Item, &user)

	// TODO: Handle groups

	return &user, nil
}

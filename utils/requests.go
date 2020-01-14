package utils

import (
	"fmt"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func getUser(id string, userData interface{}, groups []string, client dynamodb.DynamoDB) *models.User {
	user := models.User{}

	input := dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		TableName: aws.String("auth"),
	}

	result, err := client.GetItem(&input)
	if err != nil {
		fmt.Println("error", err)
		return nil
	}

	result.


	return user
}

func validateUser(user models.User) {

}

func validateRequest(req models.Request, user models.User) {

}

// InitializeRequest : Initialize the request
func InitializeRequest(req models.Request) models.User {
	var userData interface{}
	var groups []string

	if req.User.Data != nil {
		userData = req.User.Data
	}

	user := getUser(req.User.ID, userData, groups)

	validateUser(user)

	validateRequest(req, user)

	return user
}

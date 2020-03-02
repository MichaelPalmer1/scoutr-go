package aws

import (
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// GetAuth : Fetch an auth identity from the collection
// Responses:
//  - nil, nil: user does not exist
//  - nil, error: error while fetching user
//  - user, nil: found user
func (api DynamoAPI) GetAuth(id string) (*models.User, error) {
	user := &models.User{ID: id}

	// Try to find user in the auth table
	result, err := api.Client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(api.Config.AuthTable),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	})
	if err != nil {
		return nil, err
	} else if result.Item == nil {
		// Failed to find user in the table
		return nil, nil
	} else {
		// Found a user, unmarshal into user object
		err := dynamodbattribute.UnmarshalMap(result.Item, &user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// GetGroup : Fetch a group from the collection
// Responses:
//  - nil, nil: group does not exist
//  - nil, error: error while fetching group
//  - user, nil: found group
func (api DynamoAPI) GetGroup(id string) (*models.Group, error) {
	group := &models.Group{ID: id}
	result, err := api.Client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(api.Config.GroupTable),
		Key: map[string]*dynamodb.AttributeValue{
			"group_id": {S: aws.String(id)},
		},
	})
	if err != nil {
		return nil, err
	} else if result.Item == nil {
		// Group is not in the table
		return nil, nil
	} else {
		// Found group, unmarshal into group object
		err := dynamodbattribute.UnmarshalMap(result.Item, &group)
		if err != nil {
			return nil, err
		}
	}

	return group, nil
}

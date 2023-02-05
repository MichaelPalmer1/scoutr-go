package aws

import (
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetAuth : Fetch an auth identity from the collection
// Responses:
//   - nil, nil: user does not exist
//   - nil, error: error while fetching user
//   - user, nil: found user
func (api DynamoAPI) GetAuth(id string) (*types.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(api.Config.AuthTable),
		Key: map[string]dynamoTypes.AttributeValue{
			"id": &dynamoTypes.AttributeValueMemberS{
				Value: id,
			},
		},
	}

	// Try to find user in the auth table
	user, err := GetItem[types.User](api.Client, input)
	if err != nil {
		return nil, err
	} else if user == nil {
		// Failed to find user in the table
		return nil, nil
	} else {
		return user, nil
	}
}

// GetGroup : Fetch a group from the collection
// Responses:
//   - nil, nil: group does not exist
//   - nil, error: error while fetching group
//   - user, nil: found group
func (api DynamoAPI) GetGroup(id string) (*types.Group, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(api.Config.GroupTable),
		Key: map[string]dynamoTypes.AttributeValue{
			"group_id": &dynamoTypes.AttributeValueMemberS{
				Value: id,
			},
		},
	}

	group, err := GetItem[types.Group](api.Client, input)
	if err != nil {
		return nil, err
	} else if group == nil {
		// Group is not in the table
		return nil, nil
	} else {
		return group, nil
	}
}

// GetEntitlements: Fetch entitlements from the database
func (api DynamoAPI) GetEntitlements(entitlementIDs []string) ([]types.User, error) {
	// Build an IN expression that limits each expression to 100 items
	conditions := api.filtering.BuildInExpr("id", entitlementIDs, false)
	if !conditions.IsSet() {
		return nil, nil
	}

	expr, err := expression.NewBuilder().WithFilter(conditions).Build()
	if err != nil {
		return nil, err
	}

	// Build input
	input := &dynamodb.ScanInput{
		TableName:                 &api.Config.AuthTable,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	// Scan for the entitlement ids
	if entitlements, err := Scan[types.User](api.Client, input); err != nil {
		return nil, err
	} else {
		return entitlements, nil
	}
}

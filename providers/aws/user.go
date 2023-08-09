package aws

import (
	"context"

	"github.com/MichaelPalmer1/scoutr-go/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetAuth : Fetch an auth identity from the collection
// Responses:
//   - nil, nil: user does not exist
//   - nil, error: error while fetching user
//   - user, nil: found user
func (api DynamoAPI) GetAuth(id string) (*models.User, error) {
	user := &models.User{ID: id}

	// Try to find user in the auth table
	result, err := api.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(api.Config.AuthTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	})
	if err != nil {
		return nil, err
	} else if result.Item == nil {
		// Failed to find user in the table
		return nil, nil
	} else {
		// Found a user, unmarshal into user object
		err := attributevalue.UnmarshalMap(result.Item, &user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// GetGroup : Fetch a group from the collection
// Responses:
//   - nil, nil: group does not exist
//   - nil, error: error while fetching group
//   - user, nil: found group
func (api DynamoAPI) GetGroup(id string) (*models.Group, error) {
	group := &models.Group{ID: id}
	result, err := api.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(api.Config.GroupTable),
		Key: map[string]types.AttributeValue{
			"group_id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	})
	if err != nil {
		return nil, err
	} else if result.Item == nil {
		// Group is not in the table
		return nil, nil
	} else {
		// Found group, unmarshal into group object
		err := attributevalue.UnmarshalMap(result.Item, &group)
		if err != nil {
			return nil, err
		}
	}

	return group, nil
}

// GetEntitlements: Fetch entitlements from the database
func (api DynamoAPI) GetEntitlements(entitlementIDs []string) ([]models.User, error) {
	// Build an IN expression that limits each expression to 100 items
	conditions := api.Filtering.BuildInExpr("id", entitlementIDs, false)
	if conditions == nil {
		return nil, nil
	}

	conds := *conditions.(*expression.ConditionBuilder)

	expr, err := expression.NewBuilder().WithFilter(conds).Build()
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
	var entitlements []models.User
	paginator := dynamodb.NewScanPaginator(api.Client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		// Unmarshal data into User model
		records := []models.User{}
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &records); err != nil {
			return nil, err
		}

		// Append records to results
		entitlements = append(entitlements, records...)
	}

	return entitlements, nil
}

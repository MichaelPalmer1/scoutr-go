package aws_test

import (
	"context"
	"testing"

	"github.com/MichaelPalmer1/scoutr-go/pkg/config"
	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/aws"
	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cenkalti/backoff/v4"
)

type mockDynamoGetAuth struct {
	types.DynamoClientAPI
}

func (m mockDynamoGetAuth) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return &dynamodb.GetItemOutput{
		Item: map[string]dynamoTypes.AttributeValue{
			"ID": &dynamoTypes.AttributeValueMemberS{
				Value: "user-123",
			},
			"Username": &dynamoTypes.AttributeValueMemberS{
				Value: "username",
			},
		},
	}, nil
}

func TestGetAuth(t *testing.T) {
	api := aws.DynamoAPI{
		Client: mockDynamoGetAuth{},
		Scoutr: &base.Scoutr{
			Config: config.Config{
				AuthTable: "auth",
			},
		},
	}

	if user, err := api.GetAuth("user-123"); err != nil {
		t.Error(err)
	} else if user == nil {
		t.Fatal("User should not be nil")
	} else if user.ID != "user-123" {
		t.Errorf("Expected user id 'user-123' but got '%s'", user.ID)
	} else if user.Username != "username" {
		t.Errorf("Expected username 'username' but got '%s'", user.Username)
	}
}

type mockDynamoGetAuthNotFound struct {
	types.DynamoClientAPI
}

func (m mockDynamoGetAuthNotFound) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return &dynamodb.GetItemOutput{}, nil
}

func TestGetAuthNotFound(t *testing.T) {
	api := aws.DynamoAPI{
		Client: mockDynamoGetAuthNotFound{},
		Scoutr: &base.Scoutr{
			Config: config.Config{
				AuthTable: "auth",
			},
		},
	}

	if user, err := api.GetAuth("user-123"); err != nil {
		t.Error(err)
	} else if user != nil {
		t.Fatal("User should be nil")
	}
}

type mockDynamoGetAuthError struct {
	types.DynamoClientAPI
}

func (m mockDynamoGetAuthError) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return nil, backoff.Permanent(&dynamoTypes.InternalServerError{})
}

func TestGetAuthError(t *testing.T) {
	api := aws.DynamoAPI{
		Client: mockDynamoGetAuthError{},
		Scoutr: &base.Scoutr{
			Config: config.Config{
				AuthTable: "auth",
			},
		},
	}

	if _, err := api.GetAuth("user-123"); err == nil {
		t.Error("Error should not be nil")
	}
}

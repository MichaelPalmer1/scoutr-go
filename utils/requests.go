package utils

import (
	"fmt"
	"regexp"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func validateUser(user *models.User) error {
	// Make sure the user contains the required keys
	if user.ID == "" || user.Username == "" || user.Name == "" || user.Email == "" {
		return &models.Unauthorized{
			Message: "User missing one of the following fields: id, username, name, email",
		}
	}

	// TODO: Validate exclude fields

	// TODO: Validate filter fields

	// TODO: Validate permitted endpoints

	return nil
}

func canAccessEndpoint(method string, path string, user *models.User) bool {
	for _, item := range user.PermittedEndpoints {
		if method == item.Method {
			re := regexp.MustCompile(item.Endpoint)
			if re.MatchString(path) {
				return true
			}
		}
	}
	return false
}

func validateRequest(req models.Request, user *models.User) error {
	// Make sure the user has permissions to access this endpoint
	if canAccessEndpoint(req.Method, req.Path, user) {
		// TODO: Log request

		// User is authorized to access this endpoint
		return nil
	}
	// User is not authorized
	return &models.Unauthorized{
		Message: fmt.Sprintf("Not authorized to perform %s on endpoint %s", req.Method, req.Path),
	}
}

// InitializeRequest : Initialize the request
func InitializeRequest(req models.Request, client dynamodb.DynamoDB) (*models.User, error) {
	var userData *models.UserData
	var groups []string

	if req.User.Data != nil {
		userData = req.User.Data
	}

	user, err := GetUser(req.User.ID, "auth", "groups", client, userData, groups)
	if err != nil {
		return nil, err
	}

	if err := validateUser(user); err != nil {
		fmt.Println("Bad User:", err)
		return nil, err
	}

	if err := validateRequest(req, user); err != nil {
		fmt.Println("Unauthorized:", err)
		return nil, err
	}

	return user, nil
}

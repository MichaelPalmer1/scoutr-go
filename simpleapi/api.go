package simpleapi

// THIS FILE IS DEPRECATED AND MOVED TO PROVIDERS/AWS

import (
	"fmt"
	"regexp"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

// SimpleAPI : Class
type SimpleAPI struct {
	Config config.Config
	Client *dynamodb.DynamoDB
}

func scan(input *dynamodb.ScanInput, client *dynamodb.DynamoDB) ([]models.Record, error) {
	results := []models.Record{}
	err := client.ScanPages(input,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			// Unmarshal data into Record model
			records := []models.Record{}
			dynamodbattribute.UnmarshalListOfMaps(page.Items, &records)

			// Append records to results
			results = append(results, records...)

			return true
		},
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func scanAudit(input *dynamodb.ScanInput, client *dynamodb.DynamoDB) ([]models.AuditLog, error) {
	results := []models.AuditLog{}
	err := client.ScanPages(input,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			// Unmarshal data into AuditLog model
			records := []models.AuditLog{}
			dynamodbattribute.UnmarshalListOfMaps(page.Items, &records)

			// Append records to results
			results = append(results, records...)

			return true
		},
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// CanAccessEndpoint : Check whether a user has permission to access an endpoint
func (api *SimpleAPI) CanAccessEndpoint(method string, path string, user *models.User, request *models.Request) bool {
	var err error
	if request != nil {
		// Fetch the user
		user, err = utils.GetUser(request.User.ID, api.Config.AuthTable, api.Config.GroupTable, *api.Client, request.User.Data, []string{})
		if err != nil {
			log.Errorln("Failed to fetch user", err)
			return false
		}

		// Validate the user
		err = api.validateUser(user)
		if err != nil {
			log.Println("Encountered error while validating user", err)
			return false
		}
	}

	// Verify user was provided/looked up
	if user == nil {
		log.Println("Unable to validate if user has access to endpoint because user was nil")
		return false
	}

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

func (api *SimpleAPI) validateUser(user *models.User) error {
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

func (api *SimpleAPI) userIdentifier(user *models.User) string {
	return fmt.Sprintf("%s: %s (%s - %s)", user.ID, user.Name, user.Username, user.Email)
}

func (api *SimpleAPI) validateRequest(req models.Request, user *models.User) error {
	// Make sure the user has permissions to access this endpoint
	if api.CanAccessEndpoint(req.Method, req.Path, user, nil) {
		// Log request
		if req.Method == "GET" {
			log.Infof("[%s] Performed %s on %s", api.userIdentifier(user), req.Method, req.Path)
		} else {
			log.Infof("[%s] Performed %s on %s:\n%s", api.userIdentifier(user), req.Method, req.Path, req.Body)
		}

		// User is authorized to access this endpoint
		return nil
	}
	// User is not authorized
	return &models.Unauthorized{
		Message: fmt.Sprintf("Not authorized to perform %s on endpoint %s", req.Method, req.Path),
	}
}

// initializeRequest : Initialize the request
func (api *SimpleAPI) initializeRequest(req models.Request, client dynamodb.DynamoDB) (*models.User, error) {
	var userData *models.UserData
	var groups []string

	if req.User.Data != nil {
		userData = req.User.Data
	}

	user, err := utils.GetUser(req.User.ID, api.Config.AuthTable, api.Config.GroupTable, client, userData, groups)
	if err != nil {
		return nil, err
	}

	if err := api.validateUser(user); err != nil {
		log.Warnln("[%s] Bad User - %s", api.userIdentifier(user), err)
		return nil, err
	}

	if err := api.validateRequest(req, user); err != nil {
		log.Warnf("[%s] %s", api.userIdentifier(user), err)
		return nil, err
	}

	return user, nil
}

package base

import (
	"fmt"
	"regexp"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/models"
	log "github.com/sirupsen/logrus"
)

type BaseAPI interface {
	InitializeRequest(models.Request) (*models.User, error)
	GetUser(string, *models.UserData, []string) (*models.User, error)
	List(models.Request) ([]models.Record, error)
}

type SimpleAPI struct {
	BaseAPI
	Config config.Config
}

func (api *SimpleAPI) UserIdentifier(user *models.User) string {
	return fmt.Sprintf("%s: %s (%s - %s)", user.ID, user.Name, user.Username, user.Email)
}

func (api *SimpleAPI) CanAccessEndpoint(method string, path string, user *models.User, request *models.Request) bool {
	var err error
	if request != nil {
		// Fetch the user
		user, err = api.GetUser(request.User.ID, request.User.Data, []string{})
		if err != nil {
			log.Errorln("Failed to fetch user", err)
			return false
		}

		// Validate the user
		err = api.ValidateUser(user)
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

func (api *SimpleAPI) ValidateUser(user *models.User) error {
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

func (api *SimpleAPI) ValidateRequest(req models.Request, user *models.User) error {
	// Make sure the user has permissions to access this endpoint
	if api.CanAccessEndpoint(req.Method, req.Path, user, nil) {
		// Log request
		userID := api.UserIdentifier(user)
		if req.Method == "GET" {
			log.Infof("[%s] Performed %s on %s", userID, req.Method, req.Path)
		} else {
			log.Infof("[%s] Performed %s on %s:\n%s", userID, req.Method, req.Path, req.Body)
		}

		// User is authorized to access this endpoint
		return nil
	}
	// User is not authorized
	return &models.Unauthorized{
		Message: fmt.Sprintf("Not authorized to perform %s on endpoint %s", req.Method, req.Path),
	}
}

func (api *SimpleAPI) PostProcess(data []models.Record, user *models.User) {
	for _, item := range data {
		for _, key := range user.ExcludeFields {
			if _, ok := item[key]; ok {
				delete(item, key)
			}
		}
	}
}

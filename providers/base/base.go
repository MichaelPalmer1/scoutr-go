package base

import (
	"fmt"
	"regexp"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	log "github.com/sirupsen/logrus"
)

// BaseAPI : Low level interface that defines all the functions used by a SimpleAPI provider. Some of these would be
// implemented by the SimpleAPI struct
type BaseAPI interface {
	GetConfig() config.Config
	CanAccessEndpoint(string, string, *models.User, *models.Request) bool
	InitializeRequest(models.Request) (*models.User, error)
	GetUser(string, *models.UserData) (*models.User, error)
	Create(models.Request, map[string]string, map[string]utils.FieldValidation) error
	Update(models.Request, map[string]string, map[string]string, map[string]utils.FieldValidation, string) (interface{}, error)
	Get(models.Request, string) (models.Record, error)
	List(models.Request) ([]models.Record, error)
	ListUniqueValues(models.Request, string) ([]string, error)
	ListAuditLogs(models.Request, map[string]string, map[string]string) ([]models.AuditLog, error)
	History(models.Request, string, string, map[string]string, []string) ([]models.History, error)
	Search(models.Request, string, []string) ([]models.Record, error)
	Delete(models.Request, map[string]string) error
}

// SimpleAPI : Base struct that implements BaseAPI and sets up some commonly used functions across
// various cloud providers
type SimpleAPI struct {
	BaseAPI
	Filtering Filtering
	Config    config.Config
}

// GetConfig : Return config
func (api *SimpleAPI) GetConfig() config.Config {
	return api.Config
}

// UserIdentifier : Generate a user identifier for logs
func (api *SimpleAPI) UserIdentifier(user *models.User) string {
	return fmt.Sprintf("%s: %s (%s - %s)", user.ID, user.Name, user.Username, user.Email)
}

// CanAccessEndpoint : Determine if a user has access to a specific endpoint
func (api *SimpleAPI) CanAccessEndpoint(method string, path string, user *models.User, request *models.Request) bool {
	var err error
	if request != nil {
		// Fetch the user
		user, err = api.GetUser(request.User.ID, request.User.Data)
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

// ValidateUser : Validate the user object has all required fields
func (api *SimpleAPI) ValidateUser(user *models.User) error {
	// Make sure the user contains the required keys
	if user.ID == "" || user.Username == "" || user.Name == "" || user.Email == "" {
		return &models.Unauthorized{
			Message: "User missing one of the following fields: id, username, name, email",
		}
	}

	// TODO: Validate exclude fields

	// TODO: Validate filter fields

	// Make sure all the endpoints are valid regex
	for _, item := range user.PermittedEndpoints {
		if _, err := regexp.Compile(item.Endpoint); err != nil {
			return &models.BadRequest{
				Message: fmt.Sprintf("Failed to compile endpoint regex: %v", err.Error()),
			}
		}
	}

	return nil
}

// ValidateRequest : Validate the user has permissions to perform the request
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
	// User is not permitted to perform this API call
	return &models.Forbidden{
		Message: fmt.Sprintf("Not authorized to perform %s on endpoint %s", req.Method, req.Path),
	}
}

// PostProcess : Perform post processing on records before returning to user
func (api *SimpleAPI) PostProcess(data []models.Record, user *models.User) {
	for _, item := range data {
		for _, key := range user.ExcludeFields {
			if _, ok := item[key]; ok {
				delete(item, key)
			}
		}
	}
}

// MergePermissions : Merge permissions expressed in a group into the user object
func (api *SimpleAPI) MergePermissions(user *models.User, group *models.Group) {
	// Merge permitted endpoints
	user.PermittedEndpoints = append(user.PermittedEndpoints, group.PermittedEndpoints...)

	// Merge exclude fields
	user.ExcludeFields = append(user.ExcludeFields, group.ExcludeFields...)

	// Merge update fields restricted
	user.UpdateFieldsRestricted = append(user.UpdateFieldsRestricted, group.UpdateFieldsRestricted...)

	// Merge update fields permitted
	user.UpdateFieldsPermitted = append(user.UpdateFieldsPermitted, group.UpdateFieldsPermitted...)

	// Merge filter fields
	user.FilterFields = append(user.FilterFields, group.FilterFields...)
}

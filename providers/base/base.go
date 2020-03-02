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
	GetAuth(string) (*models.User, error)
	GetGroup(string) (*models.Group, error)
	CanAccessEndpoint(BaseAPI, string, string, *models.User, *models.Request) bool
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
func (api *SimpleAPI) CanAccessEndpoint(baseApi BaseAPI, method string, path string, user *models.User, request *models.Request) bool {
	var err error
	if request != nil {
		// Fetch the user
		user, err = api.GetUser(baseApi, request.User.ID, request.User.Data)
		if err != nil {
			log.Errorf("Failed to fetch user: %v", err)
			return false
		}

		// Validate the user
		err = api.ValidateUser(user)
		if err != nil {
			log.Errorf("Encountered error while validating user: %v", err)
			return false
		}
	}

	// Verify user was provided/looked up
	if user == nil {
		log.Warnln("Unable to validate if user has access to endpoint because user was nil")
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
	if api.CanAccessEndpoint(api, req.Method, req.Path, user, nil) {
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

// BuildParams : Takes in a request object and generates a parameters map
// that can be used in Filter calls
func (api *SimpleAPI) BuildParams(req models.Request) map[string]string {
	params := make(map[string]string)

	// Copy query params into params
	for key, value := range req.QueryParams {
		params[key] = value
	}

	// Merge path params into params
	for key, value := range req.PathParams {
		params[key] = value
	}

	// Generate dynamic search
	searchKey, hasSearchKey := req.PathParams["search_key"]
	searchValue, hasSearchValue := req.PathParams["search_value"]
	if hasSearchKey && hasSearchValue {
		// Map the search key and value into path params
		params[searchKey] = searchValue
		delete(params, "search_key")
		delete(params, "search_value")
	}

	return params
}

// InitializeRequest : Given a request, get the corresponding user and perform
// user and request validation.
func (api SimpleAPI) InitializeRequest(baseApi BaseAPI, req models.Request) (*models.User, error) {
	user, err := api.GetUser(baseApi, req.User.ID, req.User.Data)
	if err != nil {
		return nil, err
	}

	if err := api.ValidateUser(user); err != nil {
		log.Warnf("[%s] Bad User - %s", api.UserIdentifier(user), err)
		return nil, err
	}

	if err := api.ValidateRequest(req, user); err != nil {
		log.Warnf("[%s] %s", api.UserIdentifier(user), err)
		return nil, err
	}

	return user, nil
}

// GetUser : Fetch a user from the backend, merging any permissions from group memberships
func (api SimpleAPI) GetUser(baseApi BaseAPI, id string, userData *models.UserData) (*models.User, error) {
	isUser := true
	user := models.User{ID: id}

	// Try to find user in the auth table
	auth, err := baseApi.GetAuth(id)
	if err != nil {
		// Error while fetching user
		log.Errorf("Failed to get user: %v", err)
		return nil, err
	} else if auth == nil {
		// Failed to find user in the table
		isUser = false
	} else {
		user = *auth
	}

	// Try to find supplied entitlements in the auth table
	entitlementIDs := []string{}
	if userData != nil {
		for _, id := range userData.Groups {
			entitlement, err := baseApi.GetAuth(id)
			if err != nil {
				return nil, err
			} else if entitlement == nil {
				log.Errorln("Failed to get entitlement", err)

				// Entitlement not in the auth table
				continue
			}

			// Store this as a real entitlement
			entitlementIDs = append(entitlementIDs, id)

			// Add sub-groups
			user.Groups = append(user.Groups, entitlement.Groups...)

			// Merge permitted endpoints
			user.PermittedEndpoints = append(user.PermittedEndpoints, entitlement.PermittedEndpoints...)

			// Merge exclude fields
			user.ExcludeFields = append(user.ExcludeFields, entitlement.ExcludeFields...)

			// Merge update fields restricted
			user.UpdateFieldsRestricted = append(user.UpdateFieldsRestricted, entitlement.UpdateFieldsRestricted...)

			// Merge update fields permitted
			user.UpdateFieldsPermitted = append(user.UpdateFieldsPermitted, entitlement.UpdateFieldsPermitted...)

			// Merge filter fields
			user.FilterFields = append(user.FilterFields, entitlement.FilterFields...)
		}
	}

	// Check that a user was found
	if !isUser && len(entitlementIDs) == 0 {
		return nil, &models.Unauthorized{
			Message: fmt.Sprintf("Auth id '%s' is not authorized", id),
		}
	}

	// If the user is a member of a group, merge in the group's permissions
	for _, groupID := range user.Groups {
		group, err := baseApi.GetGroup(groupID)
		if err != nil {
			log.Errorf("Error while fetching group: %v", err)
			return nil, err
		} else if group == nil {
			// Group is not in the table
			return nil, &models.Unauthorized{
				Message: fmt.Sprintf("Group '%s' does not exist", groupID),
			}
		}

		// Merge permissions
		api.MergePermissions(&user, group)
	}

	// Save user groups before applying metadata
	userGroups := user.Groups

	// Update user object with metadata
	if userData != nil {
		if userData.Username != "" {
			user.Username = userData.Username
		}
		if userData.Name != "" {
			user.Name = userData.Name
		}
		if userData.Email != "" {
			user.Email = userData.Email
		}
		if len(userData.Groups) > 0 {
			user.Groups = userData.Groups
		}
	}

	// Update user object with all applied entitlements
	if len(entitlementIDs) > 0 {
		var groups []string
		groups = append(groups, userGroups...)
		groups = append(groups, entitlementIDs...)
		user.Groups = groups
	}

	return &user, nil
}

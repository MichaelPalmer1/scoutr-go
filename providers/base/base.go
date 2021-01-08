package base

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/MichaelPalmer1/scoutr-go/config"
	"github.com/MichaelPalmer1/scoutr-go/models"
	log "github.com/sirupsen/logrus"
)

const (
	AuditActionCreate = "CREATE"
	AuditActionUpdate = "UPDATE"
	AuditActionList   = "LIST"
	AuditActionGet    = "GET"
	AuditActionSearch = "SEARCH"
	AuditActionDelete = "DELETE"
)

// ScoutrBase : Low level interface that defines all the functions used by a Scoutr provider. Some of these would be
// implemented by the Scoutr struct
type ScoutrBase interface {
	ScoutrProvider
	GetConfig() config.Config
	CanAccessEndpoint(string, string, *models.User, *models.Request) bool
}

type ScoutrProvider interface {
	GetEntitlements([]string) ([]models.User, error)
	GetAuth(string) (*models.User, error)
	GetGroup(string) (*models.Group, error)
	Create(models.Request, map[string]interface{}, map[string]models.FieldValidation, []string) error
	Update(models.Request, map[string]string, map[string]interface{}, map[string]models.FieldValidation, []string, string) (interface{}, error)
	Get(models.Request, string) (models.Record, error)
	List(models.Request) ([]models.Record, error)
	ListUniqueValues(models.Request, string) ([]string, error)
	ListAuditLogs(models.Request, map[string]string, map[string][]string) ([]models.AuditLog, error)
	History(models.Request, string, string, map[string][]string, []string) ([]models.History, error)
	Search(models.Request, string, []string) ([]models.Record, error)
	Delete(models.Request, map[string]string) error
}

// Scoutr : Base struct that implements ScoutrBase and sets up some commonly used functions across
// various cloud providers
type Scoutr struct {
	ScoutrBase
	Filtering BaseFiltering
	Config    config.Config
}

// GetConfig : Return config
func (api *Scoutr) GetConfig() config.Config {
	return api.Config
}

// UserIdentifier : Generate a user identifier for logs
func (api *Scoutr) UserIdentifier(user *models.User) string {
	return fmt.Sprintf("%s: %s (%s - %s)", user.ID, user.Name, user.Username, user.Email)
}

// CanAccessEndpoint : Determine if a user has access to a specific endpoint
func (api *Scoutr) CanAccessEndpoint(method string, path string, user *models.User, request *models.Request) bool {
	var err error
	if request != nil {
		// Fetch the user
		user, err = api.GetUser(request.User.ID, request.User.Data)
		if err != nil {
			log.WithError(err).Error("Failed to fetch user")
			return false
		}

		// Validate the user
		if err := api.ValidateUser(user); err != nil {
			log.WithError(err).Error("Encountered error while validating user")
			return false
		}
	}

	// Verify user was provided/looked up
	if user == nil {
		log.Warnln("Unable to validate if user has access to endpoint because user was nil")
		return false
	}

	// Check permitted endpoints
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
func (api *Scoutr) ValidateUser(user *models.User) error {
	// Make sure the user contains the required keys
	if user.ID == "" || user.Username == "" || user.Name == "" || user.Email == "" {
		return &models.Unauthorized{
			Message: "User missing one of the following fields: id, username, name, email",
		}
	}

	// TODO: Validate exclude fields
	// TODO: Validate read filters
	// TODO: Validate create filters
	// TODO: Validate update filters
	// TODO: Validate delete filters

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
func (api *Scoutr) ValidateRequest(req models.Request, user *models.User) error {
	// Make sure the user has permissions to access this endpoint
	if api.CanAccessEndpoint(req.Method, req.Path, user, nil) {
		// Log request
		userID := api.UserIdentifier(user)
		if req.Method == http.MethodGet {
			log.Infof("[%s] Performed %s on %s", userID, req.Method, req.Path)
		} else {
			log.Infof("[%s] Performed %s on %s:\n%s", userID, req.Method, req.Path, req.Body)
		}

		// User is authorized to access this endpoint
		return nil
	}

	// Make sure query params have keys and values
	for key, values := range req.QueryParams {
		if key == "" {
			return &models.BadRequest{
				Message: "Query strings must have keys and values",
			}
		}

		for _, value := range values {
			if value == "" {
				return &models.BadRequest{
					Message: "Query strings must have keys and values",
				}
			}
		}
	}

	// User is not permitted to perform this API call
	return &models.Forbidden{
		Message: fmt.Sprintf("Not authorized to perform %s on endpoint %s", req.Method, req.Path),
	}
}

func (api *Scoutr) ValidateFields(validation map[string]models.FieldValidation, requiredFields []string, item map[string]interface{}, existingItem map[string]interface{}) error {
	// Check for required fields
	if len(requiredFields) > 0 {
		var missingKeys []string
		for _, key := range requiredFields {
			if _, ok := item[key]; !ok {
				missingKeys = append(missingKeys, key)
			}
		}
		if len(missingKeys) > 0 {
			return &models.BadRequest{
				Message: "Missing required fields: " + strings.Join(missingKeys, ", "),
			}
		}
	}

	// Create channels and wait group
	wg := &sync.WaitGroup{}
	ch := make(chan models.ValidationOutput, len(validation))
	done := make(chan bool, 1)

	// Trigger validation goroutines
	for key, fn := range validation {
		if _, ok := item[key]; ok {
			input := &models.ValidationInput{
				Key:          key,
				Value:        item[key],
				Item:         item,
				ExistingItem: existingItem,
			}

			// Increment wait group and start goroutine
			wg.Add(1)
			go fn(input, ch)
		}
	}

	// Wait for all validations to finish
	go func(ch chan bool) {
		wg.Wait()
		ch <- true
	}(done)

	// Create result object
	errors := make(map[string]string)

	// Receive results
	for {
		select {
		case output := <-ch:
			if output.Error != nil {
				// Validation threw an error, return the error
				wg.Done()
				return output.Error
			} else if !output.Result {
				// Validation failed, return with the error message
				errors[output.Input.Key] = output.Message
			}

			// Complete wait group item
			wg.Done()

		case <-done:
			// Return when all validations have been processed
			if len(errors) > 0 {
				return &models.BadRequest{
					Errors: errors,
				}
			}

			return nil
		}
	}
}

// PostProcess : Perform post processing on records before returning to user
func (api *Scoutr) PostProcess(data []models.Record, user *models.User) {
	for _, item := range data {
		for _, key := range user.ExcludeFields {
			if _, ok := item[key]; ok {
				delete(item, key)
			}
		}
	}
}

// MergePermissions : Merge permissions expressed in a group into the user object
func (api *Scoutr) MergePermissions(user *models.User, group *models.Group) {
	// Merge permitted endpoints
	user.PermittedEndpoints = append(user.PermittedEndpoints, group.PermittedEndpoints...)

	// Merge exclude fields
	user.ExcludeFields = append(user.ExcludeFields, group.ExcludeFields...)

	// Merge update fields restricted
	user.UpdateFieldsRestricted = append(user.UpdateFieldsRestricted, group.UpdateFieldsRestricted...)

	// Merge update fields permitted
	user.UpdateFieldsPermitted = append(user.UpdateFieldsPermitted, group.UpdateFieldsPermitted...)

	// Merge filter fields
	user.ReadFilters = append(user.ReadFilters, group.ReadFilters...)
}

// BuildParams : Takes in a request object and generates a parameters map
// that can be used in Filter calls
func (api *Scoutr) BuildParams(req models.Request) map[string][]string {
	params := make(map[string][]string)

	// Copy query params into params
	for key, values := range req.QueryParams {
		params[key] = append(params[key], values...)
	}

	// Merge path params into params
	for key, value := range req.PathParams {
		params[key] = append(params[key], value)
	}

	// Generate dynamic search
	searchKey, hasSearchKey := req.PathParams["search_key"]
	searchValue, hasSearchValue := req.PathParams["search_value"]
	if hasSearchKey && hasSearchValue {
		// Map the search key and value into path params
		params[searchKey] = []string{searchValue}
		delete(params, "search_key")
		delete(params, "search_value")
	}

	return params
}

// InitializeRequest : Given a request, get the corresponding user and perform
// user and request validation.
func (api Scoutr) InitializeRequest(req models.Request) (*models.User, error) {
	// Get user
	user, err := api.GetUser(req.User.ID, req.User.Data)
	if err != nil {
		return nil, err
	}

	// Validate user
	if err := api.ValidateUser(user); err != nil {
		log.Warnf("[%s] Bad User - %s", api.UserIdentifier(user), err)
		return nil, err
	}

	// Validate request
	if err := api.ValidateRequest(req, user); err != nil {
		log.Warnf("[%s] %s", api.UserIdentifier(user), err)
		return nil, err
	}

	return user, nil
}

func (api Scoutr) PrepareCreate(request models.Request, data map[string]interface{}, validation map[string]models.FieldValidation, requiredFields []string) (*models.User, error) {
	// Get user
	user, err := api.InitializeRequest(request)
	if err != nil {
		return nil, err
	}

	// Make sure the user has permission to update all the fields specified
	var unauthorizedFields []string
	for _, field := range user.ExcludeFields {
		for key := range data {
			if field == key {
				unauthorizedFields = append(unauthorizedFields, field)
			}
		}
	}
	if len(unauthorizedFields) > 0 {
		return nil, &models.Unauthorized{
			Message: fmt.Sprintf("Not authorized to create item with fields %+v", unauthorizedFields),
		}
	}

	// Run validation
	err = api.ValidateFields(validation, requiredFields, data, nil)
	if err != nil {
		return nil, err
	}

	// Creation filters
	localFilter := LocalFiltering{
		data: data,
	}
	results, err := localFilter.Filter(user, nil, FilterActionCreate)
	if err != nil {
		return nil, err
	}
	if results == false {
		return nil, &models.Unauthorized{
			Message: fmt.Sprintf("Unauthorized value(s) for field(s): %+v", localFilter.failedFilters),
		}
	}

	return user, nil
}

// GetUser : Fetch a user from the backend, merging any permissions from group memberships
func (api Scoutr) GetUser(id string, userData *models.UserData) (*models.User, error) {
	isUser := true
	user := models.User{ID: id}

	// Try to find user in the auth table
	if auth, err := api.ScoutrBase.GetAuth(id); err != nil {
		// Error while fetching user
		log.Errorf("Failed to get user: %v", err)
		return nil, err
	} else if auth == nil {
		// Failed to find user in the table
		isUser = false
	} else {
		user = *auth
		user.ID = id
	}

	// Try to find supplied entitlements in the auth table
	var entitlementIDs []string
	if userData != nil && len(userData.Entitlements) > 0 {
		entitlements, err := api.ScoutrBase.GetEntitlements(userData.Entitlements)
		if err != nil {
			return nil, err
		}
		for _, entitlement := range entitlements {
			// Store this as a real entitlement
			entitlementIDs = append(entitlementIDs, id)

			// Add sub-groups
			user.Groups = append(user.Groups, entitlement.Groups...)

			// Merge permissions
			user.PermittedEndpoints = append(user.PermittedEndpoints, entitlement.PermittedEndpoints...)
			user.ExcludeFields = append(user.ExcludeFields, entitlement.ExcludeFields...)
			user.UpdateFieldsRestricted = append(user.UpdateFieldsRestricted, entitlement.UpdateFieldsRestricted...)
			user.UpdateFieldsPermitted = append(user.UpdateFieldsPermitted, entitlement.UpdateFieldsPermitted...)
			user.ReadFilters = append(user.ReadFilters, entitlement.ReadFilters...)
			user.CreateFilters = append(user.CreateFilters, entitlement.CreateFilters...)
			user.UpdateFilters = append(user.UpdateFilters, entitlement.UpdateFilters...)
			user.DeleteFilters = append(user.DeleteFilters, entitlement.DeleteFilters...)
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
		group, err := api.ScoutrBase.GetGroup(groupID)
		if err != nil {
			log.WithError(err).Error("Error while fetching group")
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
		if len(userData.Entitlements) > 0 {
			user.Groups = userData.Entitlements
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

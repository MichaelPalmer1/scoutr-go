package gcp

import (
	"encoding/json"
	"fmt"

	"github.com/MichaelPalmer1/simple-api-go/models"
	log "github.com/sirupsen/logrus"
)

// InitializeRequest : Given a request, get the corresponding user and perform
// user and request validation.
func (api *FirestoreAPI) InitializeRequest(req models.Request) (*models.User, error) {
	var userData *models.UserData
	groups := []string{}

	if req.User.Data != nil {
		userData = req.User.Data
	}

	user, err := api.GetUser(req.User.ID, userData, groups)
	if err != nil {
		return nil, err
	}

	if err := api.ValidateUser(user); err != nil {
		log.Warnln("[%s] Bad User - %s", api.UserIdentifier(user), err)
		return nil, err
	}

	if err := api.ValidateRequest(req, user); err != nil {
		log.Warnf("[%s] %s", api.UserIdentifier(user), err)
		return nil, err
	}

	return user, nil
}

func (api *FirestoreAPI) GetUser(id string, userData *models.UserData, groups []string) (*models.User, error) {
	isUser := true
	user := models.User{ID: id}
	authCollection := api.client.Collection(api.Config.AuthTable)
	groupCollection := api.client.Collection(api.Config.GroupTable)

	// Try to find user in the auth table
	result, err := authCollection.Doc(id).Get(api.context)
	if err != nil {
		log.Errorln("Failed to get user", err)
		return nil, err
	} else if result.Data() == nil {
		// Failed to find user in the table
		isUser = false
	} else {
		// Found a user, unmarshal into user object
		data, err := json.Marshal(result.Data())
		if err != nil {
			log.Errorf("Failed to marshal to json: %v", err)
			return nil, err
		}
		json.Unmarshal(data, &user)
	}

	// Try to find groups in the auth table
	groupIDs := []string{}
	for _, groupID := range groups {
		var group models.User
		doc := groupCollection.Doc(groupID)
		result, err := doc.Get(api.context)
		if err != nil {
			log.Errorln("Failed to get group", err)
			return nil, err
		} else if result.Data() == nil {
			// Group is not in the table
			continue
		} else {
			// Found group, unmarshal into group object
			data, err := json.Marshal(result.Data())
			if err != nil {
				log.Errorf("Failed to marshal to json: %v", err)
				return nil, err
			}
			json.Unmarshal(data, &group)
		}

		// Store this as a real group
		groupIDs = append(groupIDs, groupID)

		// Add sub-groups
		for _, item := range group.Groups {
			user.Groups = append(user.Groups, item)
		}

		// Merge permitted endpoints
		for _, item := range group.PermittedEndpoints {
			user.PermittedEndpoints = append(user.PermittedEndpoints, item)
		}

		// Merge exclude fields
		for _, item := range group.ExcludeFields {
			user.ExcludeFields = append(user.ExcludeFields, item)
		}

		// Merge update fields restricted
		for _, item := range group.UpdateFieldsRestricted {
			user.UpdateFieldsRestricted = append(user.UpdateFieldsRestricted, item)
		}

		// Merge update fields permitted
		for _, item := range group.UpdateFieldsPermitted {
			user.UpdateFieldsPermitted = append(user.UpdateFieldsPermitted, item)
		}

		// Merge filter fields
		for _, item := range group.FilterFields {
			user.FilterFields = append(user.FilterFields, item)
		}
	}

	// Check that a user was found
	if !isUser && len(groupIDs) == 0 {
		return nil, &models.Unauthorized{
			Message: fmt.Sprintf("User '%s' is not authorized", id),
		}
	}

	// If the user is a member of a group, merge in the group's permissions
	for _, groupID := range user.Groups {
		group := models.Group{}
		result, err := groupCollection.Doc(groupID).Get(api.context)
		if err != nil {
			log.Errorln("Failed to get group", err)
			return nil, err
		} else if result.Data() == nil {
			// Group is not in the table
			return nil, &models.Unauthorized{
				Message: fmt.Sprintf("Group '%s' does not exist", groupID),
			}
		} else {
			// Found group, unmarshal into group object
			data, err := json.Marshal(result.Data())
			if err != nil {
				log.Errorf("Failed to marshal to json: %v", err)
				return nil, err
			}
			json.Unmarshal(data, &group)
		}

		// Merge permitted endpoints
		for _, item := range group.PermittedEndpoints {
			user.PermittedEndpoints = append(user.PermittedEndpoints, item)
		}

		// Merge exclude fields
		for _, item := range group.ExcludeFields {
			user.ExcludeFields = append(user.ExcludeFields, item)
		}

		// Merge update fields restricted
		for _, item := range group.UpdateFieldsRestricted {
			user.UpdateFieldsRestricted = append(user.UpdateFieldsRestricted, item)
		}

		// Merge update fields permitted
		for _, item := range group.UpdateFieldsPermitted {
			user.UpdateFieldsPermitted = append(user.UpdateFieldsPermitted, item)
		}

		// Merge filter fields
		for _, item := range group.FilterFields {
			user.FilterFields = append(user.FilterFields, item)
		}
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

	// Update user object with all applied OIDC groups
	if len(groupIDs) > 0 {
		var groups []string
		for _, groupID := range userGroups {
			groups = append(groups, groupID)
		}
		for _, groupID := range groupIDs {
			groups = append(groups, groupID)
		}
		user.Groups = groups
	}

	return &user, nil
}

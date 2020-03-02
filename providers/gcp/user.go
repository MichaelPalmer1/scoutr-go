package gcp

import (
	"encoding/json"
	"fmt"

	"github.com/MichaelPalmer1/simple-api-go/models"
	log "github.com/sirupsen/logrus"
)

// InitializeRequest : Given a request, get the corresponding user and perform
// user and request validation.
func (api FirestoreAPI) InitializeRequest(req models.Request) (*models.User, error) {
	user, err := api.GetUser(req.User.ID, req.User.Data)
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

func (api FirestoreAPI) GetUser(id string, userData *models.UserData) (*models.User, error) {
	isUser := true
	user := models.User{ID: id}
	authCollection := api.Client.Collection(api.Config.AuthTable)
	groupCollection := api.Client.Collection(api.Config.GroupTable)

	// Try to find user in the auth table
	result, err := authCollection.Doc(id).Get(api.context)
	if err != nil {
		// TODO: Need better error handling to determine if the collection exists
		log.Infof("Failed to get user: %v", err)

		// Failed to find user in the table
		isUser = false
	} else {
		// Found a user, unmarshal into user object
		data, err := json.Marshal(result.Data())
		if err != nil {
			log.Errorf("Failed to marshal to json: %v", err)
			return nil, err
		}
		err = json.Unmarshal(data, &user)
		if err != nil {
			log.Errorf("Failed to unmarshal json: %v", err)
			return nil, err
		}
	}

	// Try to find supplied entitlements in the auth table
	entitlementIDs := []string{}
	if userData != nil {
		for _, id := range userData.Groups {
			var entitlement models.User
			result, err := authCollection.Doc(id).Get(api.context)
			if err != nil {
				// TODO: Need better error handling to determine if the collection exists
				log.Errorln("Failed to get group", err)

				// User group is not in the table
				continue
			} else {
				// Found group, unmarshal into group object
				data, err := json.Marshal(result.Data())
				if err != nil {
					log.Errorf("Failed to marshal to json: %v", err)
					return nil, err
				}
				err = json.Unmarshal(data, &entitlement)
				if err != nil {
					log.Errorf("Failed to unmarshal json: %v", err)
					return nil, err
				}
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
			err = json.Unmarshal(data, &group)
			if err != nil {
				log.Errorf("Failed to unmarshal json: %v", err)
				return nil, err
			}
		}

		// Merge permissions
		api.MergePermissions(&user, &group)
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

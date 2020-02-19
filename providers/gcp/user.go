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
		log.Warnln("[%s] Bad User - %s", api.UserIdentifier(user), err)
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
		log.Errorln("Failed to get user", err)

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
				json.Unmarshal(data, &entitlement)
			}

			// Store this as a real entitement
			entitlementIDs = append(entitlementIDs, id)

			// Add sub-groups
			for _, item := range entitlement.Groups {
				user.Groups = append(user.Groups, item)
			}

			// Merge permitted endpoints
			for _, item := range entitlement.PermittedEndpoints {
				user.PermittedEndpoints = append(user.PermittedEndpoints, item)
			}

			// Merge exclude fields
			for _, item := range entitlement.ExcludeFields {
				user.ExcludeFields = append(user.ExcludeFields, item)
			}

			// Merge update fields restricted
			for _, item := range entitlement.UpdateFieldsRestricted {
				user.UpdateFieldsRestricted = append(user.UpdateFieldsRestricted, item)
			}

			// Merge update fields permitted
			for _, item := range entitlement.UpdateFieldsPermitted {
				user.UpdateFieldsPermitted = append(user.UpdateFieldsPermitted, item)
			}

			// Merge filter fields
			for _, item := range entitlement.FilterFields {
				user.FilterFields = append(user.FilterFields, item)
			}
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
			json.Unmarshal(data, &group)
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
		for _, groupID := range userGroups {
			groups = append(groups, groupID)
		}
		for _, entitlement := range entitlementIDs {
			groups = append(groups, entitlement)
		}
		user.Groups = groups
	}

	return &user, nil
}

package gcp

import (
	"encoding/json"

	"github.com/MichaelPalmer1/simple-api-go/models"
	log "github.com/sirupsen/logrus"
)

// GetAuth : Fetch an auth identity from the collection
// Responses:
//  - nil, nil: user does not exist
//  - nil, error: error while fetching user
//  - user, nil: found user
func (api FirestoreAPI) GetAuth(id string) (*models.User, error) {
	collection := api.Client.Collection(api.Config.AuthTable)
	user := &models.User{ID: id}

	// Try to find user in the auth table
	result, err := collection.Doc(id).Get(api.context)
	if err != nil {
		log.Errorf("Error while fetching user: %v", err)
		return nil, nil
	} else {
		// Found a user, unmarshal into user object
		data, err := json.Marshal(result.Data())
		if err != nil {
			log.Errorf("Failed to marshal to json: %v", err)
			return nil, err
		}
		err = json.Unmarshal(data, &user)
		if err != nil {
			log.Errorf("Failed to unmarshal from json: %v", err)
			return nil, err
		}
	}

	return user, nil
}

// GetGroup : Fetch a group from the collection
// Responses:
//  - nil, nil: group does not exist
//  - nil, error: error while fetching group
//  - user, nil: found group
func (api FirestoreAPI) GetGroup(id string) (*models.Group, error) {
	collection := api.Client.Collection(api.Config.GroupTable)
	group := &models.Group{ID: id}

	// Try to find group in the group table
	result, err := collection.Doc(id).Get(api.context)
	if err != nil {
		return nil, err
	} else {
		// Found a group, unmarshal into user object
		data, err := json.Marshal(result.Data())
		if err != nil {
			log.Errorf("Failed to marshal to json: %v", err)
			return nil, err
		}
		err = json.Unmarshal(data, &group)
		if err != nil {
			log.Errorf("Failed to unmarshal from json: %v", err)
			return nil, err
		}
	}

	return group, nil
}

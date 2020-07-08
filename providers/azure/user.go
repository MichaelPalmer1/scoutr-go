package azure

import (
	"encoding/json"

	"github.com/MichaelPalmer1/scoutr-go/models"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

// GetAuth : Fetch an auth identity from the collection
// Responses:
//  - nil, nil: user does not exist
//  - nil, error: error while fetching user
//  - user, nil: found user
func (api MongoDBAPI) GetAuth(id string) (*models.User, error) {
	collection := api.Client.C(api.Config.AuthTable)
	var user *models.User

	// Try to find user in the auth table
	var result interface{}
	if err := collection.Find(bson.M{"id": id}).One(&result); err != nil {
		if err.Error() == "not found" {
			// Failed to find user in the table
			return nil, nil
		} else {
			log.Infof("Failed to get user: %v", err)
			return nil, err
		}
	} else {
		// Found a user, unmarshal into user object
		data, err := json.Marshal(result)
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
func (api MongoDBAPI) GetGroup(id string) (*models.Group, error) {
	collection := api.Client.C(api.Config.GroupTable)
	var group *models.Group

	// Try to find group in the group table
	var result interface{}
	if err := collection.Find(bson.M{"group_id": id}).One(&result); err != nil {
		return nil, err
	} else {
		// Found a group, unmarshal into user object
		data, err := json.Marshal(result)
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

package utils

import "github.com/MichaelPalmer1/simple-api-go/models"

func getUser(id string, userData interface{}, groups []string) models.User {
	user := models.User{}

	return user
}

func validateUser(user models.User) {

}

func validateRequest(req models.Request, user models.User) {

}

// InitializeRequest : Initialize the request
func InitializeRequest(req models.Request) models.User {
	var userData interface{}
	var groups []string

	if req.User.Data != nil {
		userData = req.User.Data
	}

	user := getUser(req.User.ID, userData, groups)

	validateUser(user)

	validateRequest(req, user)

	return user
}

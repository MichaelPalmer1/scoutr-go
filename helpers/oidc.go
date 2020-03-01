package helpers

import (
	"net/http"
	"strings"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/providers/base"
)

// GetUserFromOIDC : Get user information from OIDC headers
func GetUserFromOIDC(req *http.Request, api base.BaseAPI) models.RequestUser {
	// Parse groups
	groupString := req.Header.Get(api.GetConfig().OIDCGroupHeader)
	groups := []string{}
	if groupString != "" {
		groups = strings.Split(groupString, ",")
	}

	// Generate user data
	userData := models.UserData{
		Name:     req.Header.Get(api.GetConfig().OIDCNameHeader),
		Email:    req.Header.Get(api.GetConfig().OIDCEmailHeader),
		Username: req.Header.Get(api.GetConfig().OIDCUsernameHeader),
		Groups:   groups,
	}

	return models.RequestUser{
		ID:   userData.Username,
		Data: &userData,
	}
}

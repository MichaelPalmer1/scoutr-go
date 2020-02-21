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
	groupString := req.Header.Get("Oidc-Claim-" + api.GetConfig().OIDCGroupClaim)
	groups := []string{}
	if strings.Contains(groupString, ",") {
		groups = strings.Split(groupString, ",")
	}

	// Generate user data
	userData := models.UserData{
		Name:     req.Header.Get("Oidc-Claim-" + api.GetConfig().OIDCNameClaim),
		Email:    req.Header.Get("Oidc-Claim-" + api.GetConfig().OIDCEmailClaim),
		Username: req.Header.Get("Oidc-Claim-" + api.GetConfig().OIDCUsernameClaim),
		Groups:   groups,
	}

	return models.RequestUser{
		ID:   userData.Username,
		Data: &userData,
	}
}

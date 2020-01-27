package providers

import (
	"net/http"
	"strings"
	"fmt"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/models"
)

// GetUserFromOIDC : Get user information from OIDC headers
func GetUserFromOIDC(req *http.Request, config config.Config) models.RequestUser {
	data := models.UserData{}

	// Get username
	if config.OIDCUsernameClaim != "" {
		data.Username = req.Header.Get(fmt.Sprintf("Oidc-Claim-%s", config.OIDCUsernameClaim))
	}

	// Get name
	if config.OIDCNameClaim != "" {
		data.Name = req.Header.Get(fmt.Sprintf("Oidc-Claim-%s", config.OIDCNameClaim))
	}

	// Get email
	if config.OIDCEmailClaim != "" {
		data.Email = req.Header.Get(fmt.Sprintf("Oidc-Claim-%s", config.OIDCEmailClaim))
	}

	// Check for groups
	if config.OIDCGroupClaim != "" {
		groups := req.Header.Get(fmt.Sprintf("Oidc-Claim-%s", config.OIDCGroupClaim))
		data.Groups = strings.Split(groups, ",")
	}

	// Build user object
	user := models.RequestUser{
		ID:   data.Username,
		Data: &data,
	}

	return user
}

package helpers

import (
	"net/http"
	"os"
	"strings"

	"github.com/MichaelPalmer1/scoutr-go/models"
	"github.com/MichaelPalmer1/scoutr-go/providers/base"
)

// GetUserFromOIDC : Get user information from OIDC headers
func GetUserFromOIDC(req *http.Request, api base.ScoutrBase) models.RequestUser {
	entitlements := []string{}

	// Return a dummy user when in debug mode
	if os.Getenv("DEBUG") == "true" {
		entitlementString := os.Getenv("ENTITLEMENTS")
		if entitlementString != "" {
			entitlements = strings.Split(entitlementString, ",")
		}

		return models.RequestUser{
			ID: "222222222",
			Data: &models.UserData{
				Username:     "222222222",
				Name:         "George Burdell",
				Email:        "george.p.burdell@gatech.edu",
				Entitlements: entitlements,
			},
		}
	}

	// Parse entitlements
	entitlementString := req.Header.Get(api.GetConfig().OIDCGroupHeader)
	if entitlementString != "" {
		entitlements = strings.Split(entitlementString, ",")
	}

	// Generate name
	var name []string
	for _, item := range api.GetConfig().OIDCNameHeader {
		value := req.Header.Get(item)
		if value != "" {
			name = append(name, value)
		}
	}

	// Generate user data
	userData := models.UserData{
		Name:         strings.TrimSpace(strings.Join(name, " ")),
		Email:        req.Header.Get(api.GetConfig().OIDCEmailHeader),
		Username:     req.Header.Get(api.GetConfig().OIDCUsernameHeader),
		Entitlements: entitlements,
	}

	return models.RequestUser{
		ID:   userData.Username,
		Data: &userData,
	}
}

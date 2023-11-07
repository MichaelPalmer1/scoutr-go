package helpers

import (
	"net/http"
	"os"
	"strings"

	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
)

// GetUserFromOIDC : Get user information from OIDC headers
func GetUserFromOIDC(req *http.Request, api base.ScoutrBase) types.RequestUser {
	entitlements := []string{}

	// Return a dummy user when in debug mode
	if os.Getenv("DEBUG") == "true" {
		entitlementString := os.Getenv("ENTITLEMENTS")
		if entitlementString != "" {
			entitlements = strings.Split(entitlementString, ",")
		}

		return types.RequestUser{
			ID: "222222222",
			Data: &types.UserData{
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
	userData := types.UserData{
		Name:         strings.TrimSpace(strings.Join(name, " ")),
		Email:        req.Header.Get(api.GetConfig().OIDCEmailHeader),
		Username:     req.Header.Get(api.GetConfig().OIDCUsernameHeader),
		Entitlements: entitlements,
	}

	return types.RequestUser{
		ID:   userData.Username,
		Data: &userData,
	}
}

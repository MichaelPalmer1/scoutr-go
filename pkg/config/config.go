package config

import "github.com/MichaelPalmer1/scoutr-go/pkg/types"

// Config : Various configuration
type Config struct {
	DataTable          string
	AuthTable          string
	AuditTable         string
	GroupTable         string
	PrimaryKey         string
	LogRetentionDays   int
	OIDCUsernameHeader string
	OIDCNameHeader     []string
	OIDCEmailHeader    string
	OIDCGroupHeader    string
	ErrorFunc          func(req *types.Request, user *types.User, err error)
}

// MongoConfig: Mongo-specific configuration
type MongoConfig struct {
	Config
	ConnectionString string
	Database         string
}

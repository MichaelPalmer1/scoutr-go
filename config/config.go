package config

// Config : Various configuration
type Config struct {
	DataTable          string
	AuthTable          string
	AuditTable         string
	GroupTable         string
	PrimaryKey         string
	LogRetentionDays   int
	OIDCUsernameHeader string
	OIDCNameHeader     string
	OIDCEmailHeader    string
	OIDCGroupHeader    string
}

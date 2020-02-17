package config

// Config : Various configuration
type Config struct {
	DataTable         string
	AuthTable         string
	AuditTable        string
	GroupTable        string
	PrimaryKey        string
	LogRetentionDays  int
	OIDCUsernameClaim string
	OIDCNameClaim     string
	OIDCEmailClaim    string
	OIDCGroupClaim    string
}

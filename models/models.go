package models

// Record : Data record
type Record map[string]interface{}

// RequestUser : User object
type RequestUser struct {
	ID   string
	Data *UserData
}

// Request : simple request
type Request struct {
	User        RequestUser
	Method      string
	Path        string
	Body        interface{}
	SourceIP    string
	UserAgent   string
	PathParams  map[string]string
	QueryParams map[string]string
}

// PermittedEndpoint : An endpoint
type PermittedEndpoint struct {
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
}

// FilterFields : Filter fields
type FilterFields struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type permissions struct {
	PermittedEndpoints     []PermittedEndpoint `json:"permitted_endpoints"`
	FilterFields           []FilterFields      `json:"filter_fields"`
	ExcludeFields          []string            `json:"exclude_fields"`
	UpdateFieldsPermitted  []string            `json:"update_fields_permitted"`
	UpdateFieldsRestricted []string            `json:"update_fields_restricted"`
}

// Group : Group object
type Group struct {
	ID string `json:"group_id"`
	permissions
}

// User : User object
type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Groups   []string `json:"groups"`
	permissions
}

// AuditUser : User object used in audit logs
type AuditUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	SourceIP  string `json:"source_ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// UserData : User data object
type UserData struct {
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Groups   []string `json:"groups"`
}

// AuditLog : Audit log object
type AuditLog struct {
	Time        string            `json:"time"`
	User        AuditUser         `json:"user"`
	Action      string            `json:"action"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	ExpireTime  int64             `json:"expire_time,omitempty"`
	QueryParams map[string]string `json:"query_params,omitempty"`
	Resource    map[string]string `json:"resource,omitempty"`
	Body        interface{}       `json:"body,omitempty"`
}

// History : History object
type History struct {
	Time string `json:"time"`
	Data Record `json:"data"`
}

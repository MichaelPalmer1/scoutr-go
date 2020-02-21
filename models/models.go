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
	Endpoint string `json:"endpoint" firestore:"endpoint"`
	Method   string `json:"method" firestore:"method"`
}

// FilterFields : Filter fields
type FilterFields struct {
	Field string      `json:"field" firestore:"field"`
	Value interface{} `json:"value" firestore:"value"`
}

type Permissions struct {
	PermittedEndpoints     []PermittedEndpoint `json:"permitted_endpoints" firestore:"permitted_endpoints"`
	FilterFields           []FilterFields      `json:"filter_fields" firestore:"filter_fields"`
	ExcludeFields          []string            `json:"exclude_fields" firestore:"exclude_fields"`
	UpdateFieldsPermitted  []string            `json:"update_fields_permitted" firestore:"update_fields_permitted"`
	UpdateFieldsRestricted []string            `json:"update_fields_restricted" firestore:"update_fields_restricted"`
}

// Group : Group object
type Group struct {
	ID string `json:"group_id" firestore:"group_id"`
	Permissions
}

// User : User object
type User struct {
	ID       string   `json:"id" firestore:"id"`
	Username string   `json:"username" firestore:"username"`
	Name     string   `json:"name" firestore:"name"`
	Email    string   `json:"email" firestore:"email"`
	Groups   []string `json:"groups" firestore:"groups"`
	Permissions
}

// AuditUser : User object used in audit logs
type AuditUser struct {
	ID        string `json:"id" firestore:"id"`
	Username  string `json:"username" firestore:"username"`
	Name      string `json:"name" firestore:"name"`
	Email     string `json:"email" firestore:"email"`
	SourceIP  string `json:"source_ip,omitempty" firestore:"source_ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty" firestore:"user_agent,omitempty"`
}

// UserData : User data object
type UserData struct {
	Username string   `json:"username" firestore:"username"`
	Name     string   `json:"name" firestore:"name"`
	Email    string   `json:"email" firestore:"email"`
	Groups   []string `json:"groups" firestore:"groups"`
}

// AuditLog : Audit log object
type AuditLog struct {
	Time        string            `json:"time" firestore:"time"`
	User        AuditUser         `json:"user" firestore:"user"`
	Action      string            `json:"action" firestore:"action"`
	Method      string            `json:"method" firestore:"method"`
	Path        string            `json:"path" firestore:"path"`
	ExpireTime  int64             `json:"expire_time,omitempty" firestore:"expire_time,omitempty"`
	QueryParams map[string]string `json:"query_params,omitempty" firestore:"query_params,omitempty"`
	Resource    map[string]string `json:"resource,omitempty" firestore:"resource,omitempty"`
	Body        interface{}       `json:"body,omitempty" firestore:"body,omitempty"`
}

// History : History object
type History struct {
	Time string `json:"time" firestore:"time"`
	Data Record `json:"data" firestore:"data"`
}

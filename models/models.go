package models

// Record : Data record
type Record map[string]interface{}

// RequestUserData : User Data
type RequestUserData struct {
	Groups []string
}

// RequestUser : User object
type RequestUser struct {
	ID   string `json:"id"`
	Data interface{}
}

// Request : simple request
type Request struct {
	User RequestUser
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

// Group : Group object
type Group struct {
	ID                     string              `json:"group_id"`
	PermittedEndpoints     []PermittedEndpoint `json:"permitted_endpoints"`
	FilterFields           []FilterFields      `json:"filter_fields"`
	ExcludeFields          []string            `json:"exclude_fields"`
	UpdateFieldsPermitted  []string            `json:"update_fields_permitted"`
	UpdateFieldsRestricted []string            `json:"update_fields_restricted"`
}

// User : User object
type User struct {
	ID                     string              `json:"id"`
	Groups                 []string            `json:"groups"`
	PermittedEndpoints     []PermittedEndpoint `json:"permitted_endpoints"`
	FilterFields           []FilterFields      `json:"filter_fields"`
	ExcludeFields          []string            `json:"exclude_fields"`
	UpdateFieldsPermitted  []string            `json:"update_fields_permitted"`
	UpdateFieldsRestricted []string            `json:"update_fields_restricted"`
}

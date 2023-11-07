package types

// PermittedEndpoint : An endpoint
type PermittedEndpoint struct {
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
}

// FilterField: Filter field object
type FilterField struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// Permissions: Permissions struct
type Permissions struct {
	PermittedEndpoints     []PermittedEndpoint `json:"permitted_endpoints"`
	ReadFilters            []FilterField       `json:"read_filters"`
	CreateFilters          []FilterField       `json:"create_filters"`
	UpdateFilters          []FilterField       `json:"update_filters"`
	DeleteFilters          []FilterField       `json:"delete_filters"`
	ExcludeFields          []string            `json:"exclude_fields"`
	UpdateFieldsPermitted  []string            `json:"update_fields_permitted"`
	UpdateFieldsRestricted []string            `json:"update_fields_restricted"`
}

// Group : Group object
type Group struct {
	ID string `json:"group_id"`
	Permissions
}

// User : User object
type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Groups   []string `json:"groups"`
	Permissions
}

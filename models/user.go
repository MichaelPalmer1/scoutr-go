package models

// PermittedEndpoint : An endpoint
type PermittedEndpoint struct {
	Endpoint string `json:"endpoint" firestore:"endpoint"`
	Method   string `json:"method" firestore:"method"`
}

// FilterField: Filter field object
type FilterField struct {
	Field    string      `json:"field" firestore:"field"`
	Operator string      `json:"operator" firestore:"operator"`
	Value    interface{} `json:"value" firestore:"value"`
}

// Permissions: Permissions struct
type Permissions struct {
	PermittedEndpoints     []PermittedEndpoint `json:"permitted_endpoints" firestore:"permitted_endpoints"`
	ReadFilters            []FilterField       `json:"read_filters" firestore:"read_filters"`
	CreateFilters          []FilterField       `json:"create_filters" firestore:"create_filters"`
	UpdateFilters          []FilterField       `json:"update_filters" firestore:"update_filters"`
	DeleteFilters          []FilterField       `json:"delete_filters" firestore:"delete_filters"`
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

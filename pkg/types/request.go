package types

// Record : Data record
type Record map[string]interface{}

// UserData : User data object
type UserData struct {
	Username     string   `json:"username"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Entitlements []string `json:"entitlements"`
}

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
	QueryParams map[string][]string
}

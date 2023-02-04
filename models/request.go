package models

// Record : Data record
type Record map[string]interface{}

// UserData : User data object
type UserData struct {
	Username     string   `json:"username" firestore:"username"`
	Name         string   `json:"name" firestore:"name"`
	Email        string   `json:"email" firestore:"email"`
	Entitlements []string `json:"entitlements" firestore:"entitlements"`
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

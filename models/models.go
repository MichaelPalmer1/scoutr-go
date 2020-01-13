package models

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

// User : User object
type User struct {
	ID string
}

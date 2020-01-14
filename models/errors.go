package models

// Unauthorized : User does not have permission
type Unauthorized struct {
	Message string
}

func (e *Unauthorized) Error() string {
	return e.Message
}

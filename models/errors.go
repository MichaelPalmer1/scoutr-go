package models

// Unauthorized : User does not have permission
type Unauthorized struct {
	Message string
}

func (e *Unauthorized) Error() string {
	return e.Message
}

// BadRequest : User submitted a bad request
type BadRequest struct {
	Message string
}

func (e *BadRequest) Error() string {
	return e.Message
}

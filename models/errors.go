package models

type baseError struct {
	Message string
}

// Unauthorized : User is not authenticated
type Unauthorized baseError

func (e *Unauthorized) Error() string {
	return e.Message
}

// Forbidden : User does not have permission
type Forbidden baseError

func (e *Forbidden) Error() string {
	return e.Message
}

// BadRequest : User submitted a bad request
type BadRequest baseError

func (e *BadRequest) Error() string {
	return e.Message
}

// NotFound : Item does not exist
type NotFound baseError

func (e *NotFound) Error() string {
	return e.Message
}

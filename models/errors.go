package models

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type baseError struct {
	Message string
	Errors  map[string]string
}

// Unauthorized : User is not authenticated
type Unauthorized baseError

func (e *Unauthorized) Error() string {
	if len(e.Errors) > 0 {
		bs, err := json.Marshal(e.Errors)
		if err != nil {
			log.WithError(err).Error("Failed to marshal error data")
		}

		return string(bs)
	}

	return e.Message
}

// Forbidden : User does not have permission
type Forbidden baseError

func (e *Forbidden) Error() string {
	if len(e.Errors) > 0 {
		bs, err := json.Marshal(e.Errors)
		if err != nil {
			log.WithError(err).Error("Failed to marshal error data")
		}

		return string(bs)
	}

	return e.Message
}

// BadRequest : User submitted a bad request
type BadRequest baseError

func (e *BadRequest) Error() string {
	if len(e.Errors) > 0 {
		bs, err := json.Marshal(e.Errors)
		if err != nil {
			log.WithError(err).Error("Failed to marshal error data")
		}

		return string(bs)
	}

	return e.Message
}

// NotFound : Item does not exist
type NotFound baseError

func (e *NotFound) Error() string {
	if len(e.Errors) > 0 {
		bs, err := json.Marshal(e.Errors)
		if err != nil {
			log.WithError(err).Error("Failed to marshal error data")
		}

		return string(bs)
	}

	return e.Message
}

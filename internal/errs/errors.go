package errs

import (
	"errors"
	"fmt"
)

var (
	RepositoryUnknownError = errors.New("unknown error in the repository, please check the logs")
	ServiceUnknownError    = errors.New("unknown error in the service, please check the logs")
	InvalidReferenceError  = errors.New("invalid reference")
)

type EntityNotFoundError struct {
	// Name of entity
	Name string
	// Key of entity identifier (e.g. "name", "id"). Will be used in error message.
	Key string
	// Value of entity identifier (e.g. "John", 123). Will be used in error message.
	Value any
}

func (e *EntityNotFoundError) Error() string {
	if e.Name == "" {
		e.Name = "entity"
	}

	if e.Key == "" || e.Value == nil {
		return fmt.Sprintf("%s was not found", e.Name)
	}
	return fmt.Sprintf("%s with %s %v not found", e.Name, e.Key, e.Value)
}

type EntityAlreadyExistsError struct {
	// Name of entity
	Name string
	// Key of entity identifier (e.g. "name", "id"). Will be used in error message.
	Key string
	// Value of entity identifier (e.g. "John", 123). Will be used in error message.
	Value any
}

func (e *EntityAlreadyExistsError) Error() string {
	if e.Name == "" {
		e.Name = "entity"
	}

	if e.Key == "" || e.Value == nil {
		return fmt.Sprintf("%s already exists", e.Name)
	}
	return fmt.Sprintf("%s with %s %v already exists", e.Name, e.Key, e.Value)
}

type EntityStillReferencedError struct {
	// Name of entity
	Name string
	// Key of entity identifier (e.g. "name", "id"). Will be used in error message.
	Key string
	// Value of entity identifier (e.g. "John", 123). Will be used in error message.
	Value any
}

func (e *EntityStillReferencedError) Error() string {
	if e.Name == "" {
		e.Name = "entity"
	}

	if e.Key == "" || e.Value == nil {
		return fmt.Sprintf("%s is still referenced", e.Name)
	}
	return fmt.Sprintf("%s with %s %v is still referenced", e.Name, e.Key, e.Value)
}

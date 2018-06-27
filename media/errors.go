package media

import (
	"github.com/tympanix/supper/types"
)

// ErrUnknown is an error for unrecognized media
type ErrUnknown struct{}

func (e *ErrUnknown) Error() string {
	return "media is of an unknown format"
}

// IsUnknown returns true if the error is of type *ErrUnknownMedia
func IsUnknown(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*ErrUnknown)
	return ok
}

// NewUnknownErr returns a new error indicating unknown media
func NewUnknownErr() error {
	return &ErrUnknown{}
}

// ErrExists is an error for when media already exists in a library
type ErrExists struct{}

func (e *ErrExists) Error() string {
	return "media already exists"
}

// NewExistsErr return a new error indicating a conflict with existing media
func NewExistsErr() error {
	return &ErrExists{}
}

// IsExistsErr returns true if the error is of type *ErrExists
func IsExistsErr(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*ErrExists)
	return ok
}

// Error is an error concerning some media
type Error struct {
	media types.Media
	err   error
}

// Error returns the error message of the media error
func (err *Error) Error() string {
	return err.Error()
}

// NewError returns a new erro
func NewError(err error, media types.Media) *Error {
	return &Error{
		media,
		err,
	}
}

// IsMediaError returns true if the error is of type *media.Error
func IsMediaError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*Error)
	return ok
}

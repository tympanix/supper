package media

// ErrUnknown is an error for unrecognized media
type ErrUnknown struct{}

func (e *ErrUnknown) Error() string {
	return "media is now of a known format"
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

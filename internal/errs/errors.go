package errs

import "errors"

var (
	ErrAlreadyExists = errors.New("resource already exists")
	ErrNotFound      = errors.New("resource not found")
)

type ErrServiceProblem struct {
	Err error
}

func (err ErrServiceProblem) Error() string {
	return err.Err.Error()
}

func (err ErrServiceProblem) Is(target error) bool {
	_, ok := target.(ErrServiceProblem)
	return ok
}

package user

import "fmt"

type UserError struct {
	error
}

func NewUserError(input string) *UserError {
	return &UserError{fmt.Errorf(input)}
}

func (e *UserError) Error() string {
	return e.error.Error()
}

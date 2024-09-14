package errors

import "fmt"

const (
	InternalServerError = "INTERNAL_SERVER_ERROR"

	BadRequest = "BAD_REQUEST"

	NotFound = "404"
)

type AppError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (a AppError) Error() string {
	return fmt.Sprintf("%s: %s", a.Code, a.Description)
}

func NewError(code string, description string) error {
	e := &AppError{
		Code:        code,
		Description: description,
	}
	return e
}

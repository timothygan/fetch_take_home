package errors

import "fmt"

const (
	InternalServerError = "500"

	BadRequest = "400"

	NotFound = "404"
)

type AppError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (a AppError) Error() string {
	return fmt.Sprintf("%s: %s", a.Code, a.Description)
}

func NewAppError(code string, description string) error {
	e := &AppError{
		Code:        code,
		Description: description,
	}
	return e
}

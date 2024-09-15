package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError(t *testing.T) {
	appError := AppError{
		Code:        "code",
		Description: "description",
	}
	assert.Equal(t, "code: description", appError.Error())
}

func TestNewError(t *testing.T) {
	code := "code"
	description := "description"
	appError := NewAppError(code, description)
	assert.Equal(t, "code: description", appError.Error())
}

package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDuplicateError_Error(t *testing.T) {
	err := DuplicateError{}

	assert.Equal(t, DuplicateErrorMessage, err.Error())

	_, ok := interface{}(err).(error)
	assert.True(t, ok)
}

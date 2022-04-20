package protocol

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	var err error
	err = &Error{
		Code: SystemError,
		Err:  os.ErrClosed,
	}

	e := &Error{}
	assert.True(t, errors.As(err, &e))
	assert.True(t, errors.Is(err, os.ErrClosed))
}

func TestFail(t *testing.T) {
	statusCode, body := ErrorJsonBody(SystemError)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, SystemError, body.Code)
	assert.Nil(t, body.Data)

	statusCode, body = Fail(11111, "ssss")
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, 11111, body.Code)
	assert.Equal(t, "ssss", body.Message)
}

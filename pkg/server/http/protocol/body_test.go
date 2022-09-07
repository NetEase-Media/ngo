package protocol

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	var data = "aaa"
	statusCode, body := Success(data)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, data, body.Data)
	assert.Nil(t, body.GetError())
}

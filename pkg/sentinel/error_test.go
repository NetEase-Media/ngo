package sentinel

import (
	"errors"
	"testing"

	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/stretchr/testify/assert"
)

func TestBlockError(t *testing.T) {
	userErr := errors.New("user error")
	blockErr := &BlockError{
		BlockErr: base.NewBlockErrorWithMessage(base.BlockTypeSystemFlow, "error!"),
		Err:      userErr,
	}
	assert.True(t, errors.Is(blockErr, userErr))

	var te *BlockError
	assert.True(t, errors.As(blockErr, &te))
}

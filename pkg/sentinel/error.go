package sentinel

import (
	"fmt"

	"github.com/alibaba/sentinel-golang/core/base"
)

// BlockError 用来存储sentinel的熔断错误和用户自身错误
type BlockError struct {
	BlockErr *base.BlockError
	Err      error
}

func (e *BlockError) Unwrap() error {
	return e.Err
}

func (e *BlockError) Error() string {
	return fmt.Sprintf("sentinel block error: %s, wrapped error: %v", e.BlockErr.Error(), e.Err)
}

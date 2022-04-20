package util

import "github.com/NetEase-Media/ngo/pkg/log"

// CheckError 提供简介的error判断，如果err != nil则panic
func CheckError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

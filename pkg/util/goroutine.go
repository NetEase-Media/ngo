package util

import (
	"github.com/NetEase-Media/ngo/pkg/log"
)

func GoWithRecover(handler func(), recoverHandler func(r interface{})) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(r)
				if recoverHandler != nil {
					go func() {
						defer func() {
							if p := recover(); p != nil {
								log.Error(p)
							}
						}()
						recoverHandler(r)
					}()
				}
			}
		}()
		handler()
	}()
}

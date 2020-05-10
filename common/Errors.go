package common

import "errors"

var (
	ERR_LOCK_ALREADER_REQUIRED = errors.New("锁被占用")
)

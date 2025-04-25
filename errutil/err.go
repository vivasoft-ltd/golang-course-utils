package errutil

import (
	"errors"
)

var (
	ErrEmptyRedisKeyValue = errors.New("empty redisutil key or value")
)

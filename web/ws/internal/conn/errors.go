package conn

import "errors"

var (
	ErrConnectionLost = errors.New("connection lost")
	ErrSendQueueFull  = errors.New("send queue full")
)

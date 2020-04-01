package conn

import "errors"

// 错误 string 返回
var (
	ErrConnectionLost = errors.New("connection lost")
	ErrSendQueueFull  = errors.New("send queue full")
)

package errors

const (
	ErrCustom int32 = 9999
)

var (
	ErrInvalidParam  = newError(1001, "参数错误")
	ErrInvalidToken  = newError(1002, "无效的令牌")
	ErrUserNotExists = newError(1003, "用户不存在")
	ErrPasswordError = newError(1004, "密码错误")
	ErrTokenExpired  = newError(1005, "令牌已过期")
)

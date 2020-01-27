package errors

const (
	ErrCustom int32 = 9999
)

var (
	ErrInvalidParam  = NewError(1001, "参数错误")
	ErrInvalidToken  = NewError(1002, "无效的令牌")
	ErrUserNotExists = NewError(1003, "用户不存在")
	ErrPasswordError = NewError(1004, "密码错误")
)

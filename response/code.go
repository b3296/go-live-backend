package response

const (
	CodeSuccess         = 0
	CodeInvalidParams   = 1001 // 参数错误
	CodeEmailRegistered = 1002 // 邮箱已注册
	CodeUnauthorized    = 1003 // 未授权或token失效
	CodeInternalError   = 1004 // 服务器内部错误
)

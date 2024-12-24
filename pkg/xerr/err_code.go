package xerr

// 状态码: 前三位业务，后三位功能
const (
	SERVER_COMMON_ERR   = 100001
	REQUEST_PARAM_ERROR = 100002
	DB_ERROR            = 100003
)

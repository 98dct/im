package xerr

var codeText = map[int]string{
	SERVER_COMMON_ERR:   "服务其异常，请稍后重试！",
	REQUEST_PARAM_ERROR: "请求参数有误",
	DB_ERROR:            "数据库繁忙，请稍后重试！",
}

func ErrMsg(code int) string {
	if msg, ok := codeText[code]; ok {
		return msg
	}
	return codeText[SERVER_COMMON_ERR]
}

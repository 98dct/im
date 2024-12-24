package resultx

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "",
		Data: data,
	}
}

func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
		Data: nil,
	}
}

func OkHandler(_ context.Context, v interface{}) interface{} {
	return Success(v)
}

func ErrHandler(ctx context.Context, err error) (int, any) {

	logx.WithContext(ctx).Errorf("【api】err: %v", err)
	return 200, Fail(200000, err.Error())
}

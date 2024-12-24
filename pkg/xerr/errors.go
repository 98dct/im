package xerr

import (
	"github.com/zeromicro/x/errors"
)

func New(code int, msg string) error {
	return errors.New(code, msg)
}

func NewDBErr() error {
	return New(DB_ERROR, ErrMsg(DB_ERROR))
}

func NewInternalErr() error {
	return New(SERVER_COMMON_ERR, ErrMsg(SERVER_COMMON_ERR))
}

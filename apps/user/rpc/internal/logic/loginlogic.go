package logic

import (
	"context"
	"errors"
	"im/apps/user/models"
	"im/apps/user/rpc/internal/svc"
	"im/apps/user/rpc/user"
	"im/pkg/ctxdata"
	"im/pkg/encrypt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegisted = errors.New("手机号没有注册！")
	ErrPassword         = errors.New("密码错误！")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// todo: add your logic here and delete this line
	// 1.验证用户是否注册过
	u, err := l.svcCtx.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, ErrPhoneNotRegisted
		}
		return nil, err
	}

	// 2.密码验证
	if !encrypt.ValidatePasswordHash(in.Password, u.Password.String) {
		return nil, ErrPassword
	}

	// 3.生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now,
		l.svcCtx.Config.Jwt.AccessExpire, u.Id)
	if err != nil {
		return nil, err
	}

	return &user.LoginResp{
		Id:     u.Id,
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}

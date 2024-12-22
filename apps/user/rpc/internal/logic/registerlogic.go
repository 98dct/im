package logic

import (
	"context"
	"database/sql"
	"errors"
	"im/apps/user/models"
	"im/pkg/ctxdata"
	"im/pkg/encrypt"
	"im/pkg/wuid"
	"time"

	"im/apps/user/rpc/internal/svc"
	"im/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrPhoneIsRegisted = errors.New("手机号已被注册！")

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// todo: add your logic here and delete this line

	// 1.验证用户是否注册过
	u, err := l.svcCtx.FindByPhone(l.ctx, in.Phone)
	if err != nil && err != models.ErrNotFound {
		return nil, err
	}

	if u != nil {
		return nil, ErrPhoneIsRegisted
	}

	// 2.定义用户实体
	userEntity := models.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}
	// 2.1 密码加密
	if len(in.Password) > 0 {
		passwordHash, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}

		userEntity.Password = sql.NullString{
			String: string(passwordHash),
			Valid:  true,
		}
	}

	// 3.数据入库
	_, err = l.svcCtx.Insert(l.ctx, &userEntity)
	if err != nil {
		return nil, err
	}

	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now,
		l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}

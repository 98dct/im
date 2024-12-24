package user

import (
	"context"
	"github.com/jinzhu/copier"
	"im/apps/user/rpc/user"
	"im/pkg/ctxdata"

	"im/apps/user/api/internal/svc"
	"im/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line

	uId := ctxdata.GetUId(l.ctx)

	userInfoResp, err := l.svcCtx.GetUserInfo(l.ctx, &user.GetUserInfoReq{Id: uId})
	if err != nil {
		return nil, err
	}

	var res types.User
	copier.Copy(&res, userInfoResp.User)

	return &types.UserInfoResp{Info: res}, nil
}

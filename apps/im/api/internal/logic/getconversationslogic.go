package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"im/apps/im/rpc/imclient"
	"im/pkg/ctxdata"

	"im/apps/im/api/internal/svc"
	"im/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取会话
func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// easy-chat\apps\im\api\internal\logic\getconversationslogic.go
func (l *GetConversationsLogic) GetConversations(req *types.GetConversationsReq) (resp *types.GetConversationsResp, err error) {
	// todo: add your logic here and delete this line
	uid := ctxdata.GetUId(l.ctx)

	data, err := l.svcCtx.GetConversations(l.ctx, &imclient.GetConversationsReq{
		UserId: uid,
	})

	if err != nil {
		return nil, err
	}

	var res types.GetConversationsResp
	copier.Copy(&res, &data)

	return &res, err
}

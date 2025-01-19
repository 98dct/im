package group

import (
	"context"
	"github.com/pkg/errors"
	"im/apps/im/rpc/imclient"
	"im/apps/social/rpc/socialclient"
	"im/pkg/ctxdata"

	"im/apps/social/api/internal/svc"
	"im/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创群
func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGroupLogic) CreateGroup(req *types.GroupCreateReq) (resp *types.GroupCreateResp, err error) {
	// todo: add your logic here and delete this line

	uId := ctxdata.GetUId(l.ctx)

	// 创建群
	res, err := l.svcCtx.Social.GroupCreate(l.ctx, &socialclient.GroupCreateReq{
		Name:       req.Name,
		Icon:       req.Icon,
		CreatorUid: uId,
	})
	if err != nil {
		return nil, err
	}

	if res.Id == "" {
		return nil, errors.New("群id为空！")
	}

	// 创建群会话
	_, err = l.svcCtx.Im.CreateGroupConversation(l.ctx, &imclient.CreateGroupConversationReq{
		GroupId:  res.Id,
		CreateId: uId,
	})

	return nil, err
}

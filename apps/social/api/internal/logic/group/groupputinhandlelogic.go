package group

import (
	"context"
	"errors"
	"im/apps/im/rpc/imclient"
	"im/apps/social/rpc/social"
	"im/pkg/constants"
	"im/pkg/ctxdata"

	"im/apps/social/api/internal/svc"
	"im/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请进群处理
func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(req *types.GroupPutInHandleRep) (resp *types.GroupPutInHandleResp, err error) {
	// todo: add your logic here and delete this line

	uId := ctxdata.GetUId(l.ctx)
	res, err := l.svcCtx.GroupPutInHandle(l.ctx, &social.GroupPutInHandleReq{
		GroupReqId:   req.GroupReqId,
		GroupId:      req.GroupId,
		HandleResult: req.HandleResult,
	})
	if err != nil {
		return nil, err
	}

	// 不通过群id为空！
	if res.GroupId == "" {
		return nil, errors.New("群id为空！")
	}

	// 更新群会话
	// 通过了群id不为空！
	_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, &imclient.SetUpUserConversationReq{
		SendId:   uId,
		RecvId:   res.GroupId,
		ChatType: int32(constants.GroupChatType),
	})

	return nil, err
}

package group

import (
	"context"
	"errors"
	"im/apps/im/rpc/imclient"
	"im/apps/social/rpc/socialclient"
	"im/pkg/constants"
	"im/pkg/ctxdata"

	"im/apps/social/api/internal/svc"
	"im/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请进群
func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInLogic) GroupPutIn(req *types.GroupPutInRep) (resp *types.GroupPutInResp, err error) {
	// todo: add your logic here and delete this line

	uId := ctxdata.GetUId(l.ctx)

	// 创建群
	res, err := l.svcCtx.Social.GroupPutin(l.ctx, &socialclient.GroupPutinReq{
		GroupId:    req.GroupId,
		ReqId:      uId,
		ReqMsg:     req.ReqMsg,
		ReqTime:    req.ReqTime,
		JoinSource: int32(req.JoinSource),
	})
	if err != nil {
		return nil, err
	}

	// 由于验证方式不同，所以这个地方有些情况返回群id不为空，有些情况返回群id为空
	// 群id为空需要进一步审核，验证，群id不为空，无需验证，直接加群成功
	if res.GroupId == "" {
		return nil, errors.New("群id为空！")
	}

	// 更新群会话
	// 群id不为空，则直接加群成功，无需审核处理，直接更新群会话
	_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, &imclient.SetUpUserConversationReq{
		SendId:   uId,
		RecvId:   res.GroupId,
		ChatType: int32(constants.GroupChatType),
	})

	return nil, err
}

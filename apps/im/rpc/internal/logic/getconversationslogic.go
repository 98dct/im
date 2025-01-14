package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"im/apps/im/immodels"
	"im/pkg/xerr"

	"im/apps/im/rpc/im"
	"im/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话
func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
	// todo: add your logic here and delete this line

	// 1.根据用户查询会话列表
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		if err == immodels.ErrNotFound {
			return &im.GetConversationsResp{}, nil
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find by userId err: %v, req: %v",
			err, in.UserId)
	}

	var res im.GetConversationsResp
	copier.Copy(&res, &conversations)

	// 2.根据会话列表查询具体的会话

	ids := make([]string, 0, len(conversations.ConversationList))

	for _, conversation := range conversations.ConversationList {
		ids = append(ids, conversation.ConversationId)
	}

	manyConversations, err := l.svcCtx.ConversationModel.ListByConversationIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list by conversation ids err: %v, req: %v",
			err, ids)
	}

	// 3.计算是否有未读消息
	for _, conversation := range manyConversations {
		if _, ok := res.ConversationList[conversation.ConversationId]; !ok {
			continue
		}

		if res.ConversationList[conversation.ConversationId].Total < int32(conversation.Total) {

			res.ConversationList[conversation.ConversationId].ToRead = int32(conversation.Total) - res.ConversationList[conversation.ConversationId].Total
			res.ConversationList[conversation.ConversationId].Total = int32(conversation.Total)
			res.ConversationList[conversation.ConversationId].IsShow = true
		}

	}

	return &res, nil
}

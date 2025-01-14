package logic

import (
	"context"
	"github.com/pkg/errors"
	"im/apps/im/immodels"
	"im/apps/im/rpc/im"
	"im/apps/im/rpc/internal/svc"
	"im/pkg/constants"
	"im/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新会话
func (l *PutConversationsLogic) PutConversations(in *im.PutConversationsReq) (*im.PutConversationsResp, error) {
	// todo: add your logic here and delete this line

	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find by userId err: %v, req: %v",
			err, in.UserId)
	}

	if conversations.ConversationList == nil {
		conversations.ConversationList = make(map[string]*immodels.Conversation)
	}

	for s, conversation := range in.ConversationList {
		var oldTotal int
		if conversations.ConversationList[s] != nil {
			oldTotal = conversations.ConversationList[s].Total
		}

		conversations.ConversationList[s] = &immodels.Conversation{
			ConversationId: conversation.ConversationId,
			ChatType:       constants.ChatType(conversation.ChatType),
			IsShow:         conversation.IsShow,
			Total:          int(conversation.Read) + oldTotal,
			Seq:            conversation.Seq,
		}

	}

	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "update conversations err: %v, req: %v",
			err, conversations)
	}

	return &im.PutConversationsResp{}, nil
}

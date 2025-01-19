package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"im/apps/im/immodels"
	"im/apps/im/rpc/im"
	"im/apps/im/rpc/internal/svc"
	"im/pkg/constants"
	"im/pkg/wuid"
	"im/pkg/xerr"
)

type SetUpUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 建立会话: 群聊, 私聊
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	// todo: add your logic here and delete this line

	switch constants.ChatType(in.ChatType) {
	case constants.SingleChatType:
		conversationId := wuid.CombineId(in.SendId, in.RecvId)
		conversation, err := l.svcCtx.ConversationModel.FindOne(l.ctx, conversationId)
		if err != nil {
			if err == immodels.ErrNotFound {
				err := l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
					ConversationId: conversationId,
					ChatType:       constants.ChatType(in.ChatType),
				})
				if err != nil {
					return nil, errors.Wrapf(xerr.NewDBErr(), "insert conversation err: %v, req: %v",
						err, conversationId)
				}
			} else {
				return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel FindOne err: %v, req: %v",
					err, conversationId)
			}
		} else if conversation != nil {
			return nil, nil
		}

		// 建立两者的会话
		err = l.setupUserConversation(conversationId, in.SendId, in.RecvId, constants.ChatType(in.ChatType), true)
		if err != nil {
			return nil, err
		}

		err = l.setupUserConversation(conversationId, in.RecvId, in.SendId, constants.ChatType(in.ChatType), false)
		if err != nil {
			return nil, err
		}
	case constants.GroupChatType:
		err := l.setupUserConversation(in.RecvId, in.SendId, in.RecvId, constants.ChatType(in.ChatType), true)
		if err != nil {
			return nil, err
		}
	}

	return &im.SetUpUserConversationResp{}, nil
}

func (l *SetUpUserConversationLogic) setupUserConversation(conversationId, userId, recvId string,
	chatType constants.ChatType, isShow bool) error {

	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		if err == immodels.ErrNotFound {
			conversations = &immodels.Conversations{
				ID:               primitive.NewObjectID(),
				UserId:           userId,
				ConversationList: make(map[string]*immodels.Conversation),
			}
		} else {
			return errors.Wrapf(xerr.NewDBErr(), "find by user id err: %v, req: %v", err, userId)
		}
	}

	// 更新会话记录
	if _, ok := conversations.ConversationList[conversationId]; ok {
		return nil
	}

	conversations.ConversationList[conversationId] = &immodels.Conversation{
		ConversationId: conversationId,
		ChatType:       chatType,
		IsShow:         isShow,
	}

	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "update conversations model err: %v, req: %v", err, conversations)
	}

	return nil
}

package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"im/apps/im/immodels"
	"im/apps/im/rpc/im"
	"im/apps/im/rpc/internal/svc"
	"im/pkg/constants"
	"im/pkg/xerr"
)

type CreateGroupConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupConversationLogic) CreateGroupConversation(in *im.CreateGroupConversationReq) (*im.CreateGroupConversationResp, error) {
	// todo: add your logic here and delete this line

	res := &im.CreateGroupConversationResp{}
	_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, in.GroupId)
	if err == nil {
		return res, nil
	}

	if err != immodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find one conversation err: %v, req: %v",
			err, in.GroupId)
	}

	err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
		ConversationId: in.GroupId,
		ChatType:       constants.GroupChatType,
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert one conversation err: %v, req: %v",
			err, &immodels.Conversation{
				ConversationId: in.GroupId,
				ChatType:       constants.GroupChatType,
			})
	}

	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&im.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})

	return res, err
}

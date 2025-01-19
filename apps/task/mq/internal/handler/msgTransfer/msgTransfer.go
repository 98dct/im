package msgTransfer

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"im/apps/im/ws/websocket"
	"im/apps/im/ws/ws"
	"im/apps/social/rpc/socialclient"
	"im/apps/task/mq/internal/svc"
	"im/pkg/constants"
)

type BaseMsgTransfer struct {
	svc *svc.ServiceContext
	logx.Logger
}

func NewBaseMsgTransfer(svc *svc.ServiceContext) *BaseMsgTransfer {
	return &BaseMsgTransfer{
		svc:    svc,
		Logger: logx.WithContext(context.Background()),
	}
}

func (b *BaseMsgTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	var err error
	switch data.ChatType {
	case constants.SingleChatType:
		err = b.single(ctx, data)
	case constants.GroupChatType:
		err = b.group(ctx, data)
	}

	return err
}

func (b *BaseMsgTransfer) single(ctx context.Context, data *ws.Push) error {
	return b.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (b *BaseMsgTransfer) group(ctx context.Context, data *ws.Push) error {
	// 根据群id查询所有的群用户
	users, err := b.svc.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{GroupId: data.RecvId})
	if err != nil {
		return err
	}

	data.RecvIds = make([]string, 0, len(users.List))
	for _, members := range users.List {
		if members.UserId == data.SendId {
			continue
		}
		data.RecvIds = append(data.RecvIds, members.UserId)
	}

	return b.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

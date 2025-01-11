package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"im/apps/im/immodels"
	"im/apps/im/ws/websocket"
	"im/apps/task/mq/internal/svc"
	"im/apps/task/mq/mq"
	"im/pkg/constants"
)

type MsgChatTransfer struct {
	logx.Logger
	svc *svc.ServiceContext
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svc:    svc,
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key: ", key, "value: ", value)

	var (
		data mq.MsgChatTransfer
	)
	err := json.Unmarshal([]byte(value), &data)
	if err != nil {
		return err
	}

	// 存储聊天信息
	if err := m.addChatLog(ctx, data); err != nil {
		return err
	}

	// 推送消息
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})

}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, data mq.MsgChatTransfer) error {
	// 记录消息
	msg := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		MsgFrom:        0,
		ChatType:       data.ChatType,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}

	return m.svc.ChatLogModel.Insert(ctx, &msg)
}

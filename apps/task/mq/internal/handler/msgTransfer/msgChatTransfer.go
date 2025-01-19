package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"im/apps/im/immodels"
	"im/apps/im/ws/ws"
	"im/apps/task/mq/internal/svc"
	"im/apps/task/mq/mq"
	"im/pkg/bitmap"
)

type MsgChatTransfer struct {
	*BaseMsgTransfer
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svc),
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key: ", key, "value: ", value)

	var (
		data  mq.MsgChatTransfer
		msgId = primitive.NewObjectID()
	)
	err := json.Unmarshal([]byte(value), &data)
	if err != nil {
		return err
	}

	// 存储聊天信息
	if err := m.addChatLog(ctx, msgId, data); err != nil {
		return err
	}

	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MType:          data.MType,
		MsgId:          msgId.Hex(),
		Content:        data.Content,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data mq.MsgChatTransfer) error {
	// 记录消息
	msg := immodels.ChatLog{
		ID:             msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		MsgFrom:        0,
		ChatType:       data.ChatType,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}

	err := m.svc.ChatLogModel.Insert(ctx, &msg)
	if err != nil {
		return err
	}

	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(msg.SendId)
	msg.ReadRecords = readRecords.Export()

	// 更新会话
	return m.svc.ConversationModel.UpdateMsg(ctx, &msg)
}

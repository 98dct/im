package mqclient

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"im/apps/task/mq/mq"
)

type MsgChatTransferClient interface {
	Push(msg *mq.MsgChatTransfer) error
}

type msgChatTransferClient struct {
	pusher *kq.Pusher
}

func NewMsgChatTransferClient(addr []string, topic string, opts ...kq.PushOption) *msgChatTransferClient {
	return &msgChatTransferClient{pusher: kq.NewPusher(addr, topic, opts...)}
}

func (m *msgChatTransferClient) Push(msg *mq.MsgChatTransfer) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return m.pusher.Push(context.Background(), string(bytes))
}

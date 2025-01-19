package msgTransfer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"im/apps/im/ws/ws"
	"im/apps/task/mq/internal/svc"
	"im/apps/task/mq/mq"
	"im/pkg/bitmap"
	"im/pkg/constants"
	"sync"
	"time"
)

var (
	GroupMsgReadRecordDelayTime  = time.Second
	GroupMsgReadRecordDelayCount = 10
)

const (
	GroupMsgReadHandlerAtTransfer = iota
	GroupMsgReadHandlerDelayTransfer
)

type MsgReadTransfer struct {
	*BaseMsgTransfer
	cache.Cache

	mu sync.Mutex

	groupMsgs map[string]*groupMsgRead
	push      chan *ws.Push
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	m := &MsgReadTransfer{
		BaseMsgTransfer: NewBaseMsgTransfer(svc),
		groupMsgs:       make(map[string]*groupMsgRead, 1),
		push:            make(chan *ws.Push, 1),
	}

	if svc.Config.MsgReadHandler.GroupMsgReadHandler != GroupMsgReadHandlerAtTransfer {
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount > 0 {
			GroupMsgReadRecordDelayCount = svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount
		}

		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
			GroupMsgReadRecordDelayTime = time.Duration(svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime) * time.Second
		}
	}

	go m.transfer()

	return m
}

func (m *MsgReadTransfer) Consume(ctx context.Context, key, value string) error {
	m.Info("msgReadTransfer ", value)
	var (
		data mq.MsgMarkRead
	)
	err := json.Unmarshal([]byte(value), &data)
	if err != nil {
		fmt.Println("err1:", err)
		return err
	}

	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		fmt.Println("err2:", err)
		return err
	}

	fmt.Println("MsgReadTransfer: ", readRecords)

	push := &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMakeRead,
		ReadRecords:    readRecords,
	}

	switch push.ChatType {
	case constants.SingleChatType:
		// 直接推送
		m.push <- push
	case constants.GroupChatType:
		// 判断是否开启合并消息的处理
		if m.svc.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			m.push <- push
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		push.SendId = ""

		if _, ok := m.groupMsgs[push.ConversationId]; ok {
			m.Infof("merge push %v", push.ConversationId)
			// 合并请求
			m.groupMsgs[push.ConversationId].merge(push)
		} else {
			m.Infof("newGroupMsgRead push %v", push.ConversationId)
			m.groupMsgs[push.ConversationId] = NewGroupMsgRead(push, m.push)
		}
	}

	return nil
}

func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	res := make(map[string]string)

	chatLogs, err := m.BaseMsgTransfer.svc.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return res, err
	}

	for _, chatLog := range chatLogs {
		switch data.ChatType {
		case constants.SingleChatType:
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			readRecords := bitmap.Load(chatLog.ReadRecords)
			readRecords.Set(data.SendId)
			chatLog.ReadRecords = readRecords.Export()
		}

		res[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)
		err := m.svc.ChatLogModel.UpdateMakeRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			m.Errorf("udpate msg read err: %v", err)
			//return res, err
		}
	}

	return res, nil

}

func (m *MsgReadTransfer) transfer() {
	for push := range m.push {
		if push.RecvId != "" || len(push.RecvIds) > 0 {
			if err := m.Transfer(context.Background(), push); err != nil {
				m.Errorf("m transfer err %v push %v", err, push)
			}
		}

		if push.ChatType == constants.SingleChatType {
			continue
		}

		if m.svc.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			continue
		}
		// 清空数据
		m.mu.Lock()
		//
		if _, ok := m.groupMsgs[push.ConversationId]; ok && m.groupMsgs[push.ConversationId].IsIdle() {
			m.groupMsgs[push.ConversationId].clear()
			delete(m.groupMsgs, push.ConversationId)
		}

		m.mu.Unlock()

	}
}

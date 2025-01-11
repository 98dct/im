package mq

import "im/pkg/constants"

type MsgChatTransfer struct {
	ConversationId     string `json:"conversationId"`
	constants.ChatType `json:"chatType"`
	SendId             string `json:"sendId"`
	RecvId             string `json:"RecvId"`
	SendTime           int64  `json:"sendTime"`
	constants.MType    `json:"mType"`
	Content            string `json:"content"`
}

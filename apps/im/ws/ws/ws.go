package ws

import (
	"im/pkg/constants"
)

type Msg struct {
	MsgId           string            `mapstructure:"msgId"`
	ReadRecords     map[string]string `mapstructure:"readRecords"`
	constants.MType `mapstructure:"mType"`
	Content         string `mapstructure:"content"`
}

type Chat struct {
	ConversationId     string `mapstructure:"conversationId"`
	constants.ChatType `mapstructure:"chatType"`
	SendId             string `mapstructure:"sendId"`
	RecvId             string `mapstructure:"recvId"`
	Msg                `mapstructure:"msg"`
	SendTime           int64 `mapstructure:"sendTime"`
}

type Push struct {
	ConversationId     string `mapstructure:"conversationId"`
	constants.ChatType `mapstructure:"chatType"`

	ReadRecords map[string]string     `mapstructure:"readRecords"`
	ContentType constants.ContentType `mapstructure:"contentType"`

	MsgId string `mapstructure:"msgId"`

	SendId          string   `mapstructure:"sendId"`
	RecvId          string   `mapstructure:"recvId"`
	RecvIds         []string `mapstructure:"recvIds"`
	SendTime        int64    `mapstructure:"sendTime"`
	constants.MType `mapstructure:"mType"`
	Content         string `mapstructure:"content"`
}

type MarkRead struct {
	constants.ChatType `mapstructure:"chatType"`
	RecvId             string   `mapstructure:"recvId"`
	ConversationId     string   `mapstructure:"conversationId"`
	MsgIds             []string `mapstructure:"msgIds"`
}

package svc

import (
	"im/apps/im/immodels"
	"im/apps/im/ws/internal/config"
	"im/apps/task/mq/mqclient"
)

type ServiceContext struct {
	Config config.Config
	immodels.ChatLogModel
	mqclient.MsgChatTransferClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{Config: c,
		ChatLogModel:          immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
	}
}

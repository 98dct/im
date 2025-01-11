package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"im/apps/im/immodels"
	websocketx "im/apps/im/ws/websocket"
	"im/apps/task/mq/internal/config"
	"im/pkg/constants"
	"net/http"
)

type ServiceContext struct {
	config.Config

	*redis.Redis
	immodels.ChatLogModel
	WsClient websocketx.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:       c,
		Redis:        redis.MustNewRedis(c.Redisx),
		ChatLogModel: immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}

	token, err := svc.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
	if err != nil {
		panic(err)
	}
	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = websocketx.NewClient(c.Ws.Host, websocketx.WithHeader(header))

	return svc
}

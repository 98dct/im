package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"im/apps/im/immodels"
	websocketx "im/apps/im/ws/websocket"
	"im/apps/social/rpc/socialclient"
	"im/apps/task/mq/internal/config"
	"im/pkg/constants"
	"net/http"
)

type ServiceContext struct {
	config.Config

	*redis.Redis
	immodels.ChatLogModel
	immodels.ConversationModel
	WsClient websocketx.Client
	socialclient.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ChatLogModel:      immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel: immodels.MustConversationModel(c.Mongo.Url, c.Mongo.Db),
		Social:            socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}

	token, err := svc.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
	if err != nil {
		panic(err)
	}
	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = websocketx.NewClient(c.Ws.Host,
		websocketx.WithHeader(header),
		websocketx.WithDiscover(websocketx.NewRedisDiscover(header, constants.REDIS_DISCOVER_SRV, c.Redisx)),
	)

	return svc
}

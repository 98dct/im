package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"im/apps/im/rpc/imclient"
	"im/apps/social/api/internal/config"
	"im/apps/social/rpc/socialclient"
	"im/apps/user/rpc/userclient"
)

type ServiceContext struct {
	Config config.Config
	socialclient.Social
	userclient.User
	imclient.Im
	*redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Im:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		Redis:  redis.MustNewRedis(c.Redisx),
	}
}

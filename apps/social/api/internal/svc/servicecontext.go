package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"im/apps/social/api/internal/config"
	"im/apps/social/rpc/socialclient"
	"im/apps/user/rpc/userclient"
)

type ServiceContext struct {
	Config config.Config
	socialclient.Social
	userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}

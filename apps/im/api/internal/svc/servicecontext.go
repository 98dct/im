package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"im/apps/im/api/internal/config"
	"im/apps/im/rpc/imclient"
	"im/apps/social/rpc/socialclient"
	"im/apps/user/rpc/userclient"
)

type ServiceContext struct {
	Config config.Config
	imclient.Im
	socialclient.Social
	userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Im:     imclient.NewIm(zrpc.MustNewClient(c.IMRpc)),
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}

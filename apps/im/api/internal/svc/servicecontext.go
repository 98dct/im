package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"im/apps/im/api/internal/config"
	"im/apps/im/rpc/imclient"
)

type ServiceContext struct {
	Config config.Config
	imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Im:     imclient.NewIm(zrpc.MustNewClient(c.IMRpc)),
	}
}

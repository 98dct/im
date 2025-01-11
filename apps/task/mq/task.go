package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"im/apps/task/mq/internal/config"
	"im/apps/task/mq/internal/handler"
	"im/apps/task/mq/internal/svc"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)

	listener := handler.NewListener(ctx)
	serviceGroup := service.NewServiceGroup()
	for _, s := range listener.Services() {
		serviceGroup.Add(s)
	}

	fmt.Println("queue already start ... ")
	defer serviceGroup.Stop()
	serviceGroup.Start()

}

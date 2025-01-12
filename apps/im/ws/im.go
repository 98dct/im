package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"im/apps/im/ws/internal/config"
	"im/apps/im/ws/internal/handler"
	"im/apps/im/ws/internal/svc"
	"im/apps/im/ws/websocket"
)

var configFile = flag.String("f", "etc/dev/ws.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
		websocket.WithServerAck(websocket.RigorAck),
		//websocket.WithServerMaxIdleTime(10*time.Second),
	)
	defer srv.Stop()

	handler.RegisterRoutes(srv, ctx)

	fmt.Println("websocket start on ", c.ListenOn)
	srv.Start()
}

package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"im/apps/im/ws/internal/config"
	"im/apps/im/ws/internal/handler"
	"im/apps/im/ws/internal/svc"
	"im/apps/im/ws/websocket"
	"im/pkg/constants"
	"im/pkg/ctxdata"
	"net/http"
	"time"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)
	// 设置服务认证的token
	token, err := ctxdata.GetJwtToken(c.JwtAuth.AccessSecret, time.Now().Unix(), 3153600000, fmt.Sprintf("ws:%s", time.Now().Unix()))
	if err != nil {
		panic(err)
	}
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
		websocket.WithServerAck(websocket.NoAck),
		websocket.WithServerDiscover(websocket.NewRedisDiscover(http.Header{
			"Authorization": []string{token},
		}, constants.REDIS_DISCOVER_SRV, c.Redisx)),
		//websocket.WithServerMaxIdleTime(10*time.Second),
	)
	defer srv.Stop()

	handler.RegisterRoutes(srv, ctx)

	fmt.Println("websocket start on ", c.ListenOn)
	srv.Start()
}

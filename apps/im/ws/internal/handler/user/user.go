package user

import (
	"im/apps/im/ws/internal/svc"
	websocketx "im/apps/im/ws/websocket"
)

func Online(svc *svc.ServiceContext) websocketx.HandlerFunc {
	return func(srv *websocketx.Server, conn *websocketx.Conn, msg *websocketx.Message) {
		users := srv.GetUsers()
		uId := srv.GetUser(conn)
		err := srv.Send(websocketx.NewMsg(uId, users), conn)
		srv.Info("err", err)
	}
}

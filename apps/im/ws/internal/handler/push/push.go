package push

import (
	"github.com/mitchellh/mapstructure"
	"im/apps/im/ws/internal/svc"
	"im/apps/im/ws/websocket"
	"im/apps/im/ws/ws"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {

		var data ws.Push
		err := mapstructure.Decode(msg.Data, &data)
		if err != nil {
			srv.Send(websocket.NewErrMsg(err), srv.GetConn(data.SendId))
			return
		}

		rconn := srv.GetConn(data.RecvId)
		if rconn == nil {
			// todo 目标离线
			return
		}

		srv.Infof("push message %v", data)

		srv.Send(websocket.NewMsg(data.SendId, ws.Chat{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			Msg: ws.Msg{
				MType:   data.MType,
				Content: data.Content,
			},
			SendTime: data.SendTime,
		}), rconn)

	}
}

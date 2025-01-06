package conversation

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"im/apps/im/ws/internal/logic"
	"im/apps/im/ws/internal/svc"
	"im/apps/im/ws/websocket"
	"im/apps/im/ws/ws"
	"im/pkg/constants"
	"time"
)

func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {

		var data ws.Chat
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMsg(err), conn)
			return
		}

		switch data.ChatType {
		case constants.SingleCHatType:
			err := logic.NewConversation(context.Background(), srv, svc).SingleChat(&data, conn.Uid)
			if err != nil {
				srv.Send(websocket.NewErrMsg(err), conn)
				return
			}

			srv.SendByUserId(websocket.NewMsg(conn.Uid, ws.Chat{
				ConversationId: data.ConversationId,
				ChatType:       data.ChatType,
				SendId:         conn.Uid,
				RecvId:         data.RecvId,
				Msg:            data.Msg,
				SendTime:       time.Now().UnixNano(),
			}), data.RecvId)
		case constants.GroupChatType:

		}

	}
}

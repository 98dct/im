package push

import (
	"github.com/mitchellh/mapstructure"
	"im/apps/im/ws/internal/svc"
	"im/apps/im/ws/websocket"
	"im/apps/im/ws/ws"
	"im/pkg/constants"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {

		var data ws.Push
		err := mapstructure.Decode(msg.Data, &data)
		if err != nil {
			srv.Send(websocket.NewErrMsg(err), srv.GetConn(data.SendId))
			return
		}

		switch data.ChatType {
		case constants.SingleChatType:
			single(srv, &data, data.RecvId)
		case constants.GroupChatType:
			group(srv, &data)
		}

	}
}

func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	rconn := srv.GetConn(recvId)
	if rconn == nil {
		// todo 目标离线
		return nil
	}

	srv.Infof("push message %v", data)

	return srv.Send(websocket.NewMsg(data.SendId, ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		Msg: ws.Msg{
			MsgId:       data.MsgId,
			ReadRecords: data.ReadRecords,
			MType:       data.MType,
			Content:     data.Content,
		},
		SendTime: data.SendTime,
	}), rconn)
}

func group(srv *websocket.Server, data *ws.Push) error {

	for _, recvId := range data.RecvIds {
		recvId := recvId
		srv.Schedule(func() {
			single(srv, data, recvId)
		})
	}
	return nil
}

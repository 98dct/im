package handler

import (
	"im/apps/im/ws/internal/handler/conversation"
	"im/apps/im/ws/internal/handler/push"
	"im/apps/im/ws/internal/handler/user"
	"im/apps/im/ws/internal/svc"
	"im/apps/im/ws/websocket"
)

func RegisterRoutes(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Router{
		{
			Method:  "user.online",
			Handler: user.Online(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svc),
		},
		{
			Method:  "push",
			Handler: push.Push(svc),
		},
	})
}

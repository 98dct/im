package handler

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"im/apps/task/mq/internal/handler/msgTransfer"
	"im/apps/task/mq/internal/svc"
)

type Listener struct {
	svc *svc.ServiceContext
}

func NewListener(svc *svc.ServiceContext) *Listener {
	return &Listener{svc: svc}
}

func (l *Listener) Services() []service.Service {
	return []service.Service{
		// todo 此处可以加载多个消费者
		kq.MustNewQueue(l.svc.Config.MsgChatTransfer, msgTransfer.NewMsgChatTransfer(l.svc)),
	}
}

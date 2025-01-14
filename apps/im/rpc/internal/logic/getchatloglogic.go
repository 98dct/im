package logic

import (
	"context"
	"github.com/pkg/errors"
	"im/pkg/xerr"

	"im/apps/im/rpc/im"
	"im/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话记录
func (l *GetChatLogLogic) GetChatLog(in *im.GetChatLogReq) (*im.GetChatLogResp, error) {
	// todo: add your logic here and delete this line

	// 根据msgid查询
	if in.MsgId != "" {
		chatLog, err := l.svcCtx.ChatLogModel.FindOne(l.ctx, in.MsgId)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "find chatlog by msgId err: %v, req: %v",
				err, in.MsgId)
		}
		return &im.GetChatLogResp{List: []*im.ChatLog{
			{
				Id:             chatLog.ID.Hex(),
				ConversationId: chatLog.ConversationId,
				SendId:         chatLog.SendId,
				RecvId:         chatLog.RecvId,
				MsgType:        int32(chatLog.MsgType),
				MsgContent:     chatLog.MsgContent,
				ChatType:       int32(chatLog.ChatType),
				SendTime:       chatLog.SendTime,
			},
		}}, nil
	}

	// 时间段分段查询
	chatLogs, err := l.svcCtx.ChatLogModel.ListBySendTime(l.ctx, in.ConversationId, in.StartSendTime, in.EndSendTime, in.Count)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list chatlogs by send time err: %v, req: %v",
			err, in)
	}

	resp := make([]*im.ChatLog, 0, len(chatLogs))
	for _, chatLog := range chatLogs {
		resp = append(resp, &im.ChatLog{
			Id:             chatLog.ID.Hex(),
			ConversationId: chatLog.ConversationId,
			SendId:         chatLog.SendId,
			RecvId:         chatLog.RecvId,
			MsgType:        int32(chatLog.MsgType),
			MsgContent:     chatLog.MsgContent,
			ChatType:       int32(chatLog.ChatType),
			SendTime:       chatLog.SendTime,
		})
	}
	return &im.GetChatLogResp{List: resp}, nil
}

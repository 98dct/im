package friend

import (
	"context"
	"im/apps/social/rpc/socialclient"
	"im/pkg/ctxdata"
	"time"

	"im/apps/social/api/internal/svc"
	"im/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请
func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInLogic) FriendPutIn(req *types.FriendPutInReq) (resp *types.FriendPutInResp, err error) {
	// todo: add your logic here and delete this line

	uId := ctxdata.GetUId(l.ctx)

	_, err = l.svcCtx.FriendPutIn(l.ctx, &socialclient.FriendPutInReq{
		UserId:  uId,
		ReqUid:  req.UserId,
		ReqMsg:  req.ReqMsg,
		ReqTime: time.Now().Unix(),
	})

	return
}

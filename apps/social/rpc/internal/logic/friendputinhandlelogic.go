package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"im/apps/social/socialmodels"
	"im/pkg/constants"
	"im/pkg/xerr"
	"time"

	"im/apps/social/rpc/internal/svc"
	"im/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrAlreadyPass = xerr.NewMsg("好友申请已经通过")
var ErrAlreadyReject = xerr.NewMsg("好友申请已经拒绝")

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// todo: add your logic here and delete this line

	// 1.获取好友申请记录
	friendRequests, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, uint64(in.FriendReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend request by friend req id fail err: %v, req: %v",
			err, in)
	}

	// 2.验证是否有处理
	switch constants.HandlerResult(friendRequests.HandleResult.Int64) {
	case constants.Pass:
		return nil, ErrAlreadyPass
	case constants.Reject:
		return nil, ErrAlreadyReject
	}
	// 3.修改处理结果----插入朋友两条记录  事务

	friendRequests.HandleResult.Int64 = int64(in.HandleResult)
	friendRequests.HandledAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 3.1 更新处理结果
		if err := l.svcCtx.FriendRequestsModel.Update(l.ctx, session, friendRequests); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend request err %v, req %v", err, friendRequests)
		}

		if in.HandleResult != int32(constants.Pass) {
			return nil
		}

		friends := []*socialmodels.Friends{
			{
				UserId:    friendRequests.UserId,
				FriendUid: friendRequests.ReqUid,
			},
			{
				UserId:    friendRequests.ReqUid,
				FriendUid: friendRequests.UserId,
			},
		}

		// 3.2 插入两条记录
		_, err := l.svcCtx.FriendsModel.Inserts(ctx, session, friends...)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert friends  err %v, req %v", err, friends)
		}
		return nil
	})

	return &social.FriendPutInHandleResp{}, nil
}

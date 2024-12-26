package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"im/apps/social/socialmodels"
	"im/pkg/constants"
	"im/pkg/xerr"
	"time"

	"im/apps/social/rpc/internal/svc"
	"im/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// todo: add your logic here and delete this line

	// 1.申请人是否与目标是好友关系
	friends, err := l.svcCtx.FindByUidAndFid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend by uid fid err %v req %v", err, in)
	}

	if friends != nil {
		return &social.FriendPutInResp{}, errors.New("已存在好友关系！")
	}
	// 2.是否有在途申请、或者之前申请不成功
	friendRequests, err := l.svcCtx.FriendRequestsModel.FindByUserIdAndReqId(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend requests by uid requid err %v req %v", err, in)
	}

	if friendRequests != nil {
		return &social.FriendPutInResp{}, errors.New("已提交好友申请！")
	}

	// todo 之前申请不成功呢

	// 3.创建申请记录
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &socialmodels.FriendRequests{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			Valid:  true,
			String: in.ReqMsg,
		},
		ReqTime: time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{
			Valid: true,
			Int64: int64(constants.NoHandler),
		},
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friend requests  err %v req %v", err, in)
	}

	return &social.FriendPutInResp{}, nil
}

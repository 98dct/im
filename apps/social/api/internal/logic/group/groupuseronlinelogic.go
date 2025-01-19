package group

import (
	"context"
	"github.com/pkg/errors"
	"im/apps/social/rpc/socialclient"
	"im/pkg/constants"
	"im/pkg/xerr"

	"im/apps/social/api/internal/svc"
	"im/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群用户在线情况
func NewGroupUserOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserOnlineLogic {
	return &GroupUserOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserOnlineLogic) GroupUserOnline(req *types.GroupUserOnlineReq) (resp *types.GroupUserOnlineResp, err error) {
	// todo: add your logic here and delete this line

	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
		GroupId: req.GroupId,
	})

	// 还需要获取用户的信息
	uids := make([]string, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {
		uids = append(uids, v.UserId)
	}

	onlines, err := l.svcCtx.Redis.Hgetall(constants.REDIS_ONLINE_USER)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "redis hgetall %v err %v", constants.REDIS_ONLINE_USER, err)
	}

	resOnLineList := make(map[string]bool, len(uids))
	for _, s := range uids {
		if _, ok := onlines[s]; ok {
			resOnLineList[s] = true
		} else {
			resOnLineList[s] = false
		}
	}

	return &types.GroupUserOnlineResp{
		OnlineList: resOnLineList,
	}, nil
}

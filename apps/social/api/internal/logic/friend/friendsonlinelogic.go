package friend

import (
	"context"
	"github.com/pkg/errors"
	"im/apps/social/rpc/socialclient"
	"im/pkg/constants"
	"im/pkg/ctxdata"
	"im/pkg/xerr"

	"im/apps/social/api/internal/svc"
	"im/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendsOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友在线情况
func NewFriendsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendsOnlineLogic {
	return &FriendsOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendsOnlineLogic) FriendsOnline(req *types.FriendsOnlineReq) (resp *types.FriendsOnlineResp, err error) {
	// todo: add your logic here and delete this line
	uid := ctxdata.GetUId(l.ctx)

	friendList, err := l.svcCtx.Social.FriendList(l.ctx, &socialclient.FriendListReq{UserId: uid})
	if err != nil {
		return nil, err
	}

	if len(friendList.List) == 0 {
		return &types.FriendsOnlineResp{}, nil
	}

	// 还需要获取用户的信息
	uids := make([]string, 0, len(friendList.List))
	for _, friends := range friendList.List {
		uids = append(uids, friends.FriendUid)
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

	return &types.FriendsOnlineResp{OnlineList: resOnLineList}, nil
}

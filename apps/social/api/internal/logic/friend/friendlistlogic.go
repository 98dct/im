package friend

import (
	"context"
	"im/apps/social/rpc/socialclient"
	"im/apps/user/rpc/userclient"
	"im/pkg/ctxdata"

	"im/apps/social/api/internal/svc"
	"im/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友列表
func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	// todo: add your logic here and delete this line

	uid := ctxdata.GetUId(l.ctx)
	friendList, err := l.svcCtx.Social.FriendList(l.ctx, &socialclient.FriendListReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	if len(friendList.List) == 0 {
		return &types.FriendListResp{}, nil
	}

	uIds := make([]string, 0, len(friendList.List))
	for _, friends := range friendList.List {
		uIds = append(uIds, friends.FriendUid)
	}

	users, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{Ids: uIds})
	if err != nil {
		return nil, err
	}

	usersMp := make(map[string]*userclient.UserEntity)
	for _, userEntity := range users.User {
		usersMp[userEntity.Id] = userEntity
	}

	friends := make([]*types.Friends, 0, len(uIds))
	for i := 0; i < len(uIds); i++ {
		friends = append(friends, &types.Friends{
			Id:        friendList.List[i].Id,
			FriendUid: friendList.List[i].FriendUid,
			Nickname:  usersMp[friendList.List[i].FriendUid].Nickname,
			Avatar:    usersMp[friendList.List[i].FriendUid].Avatar,
		})
	}

	return &types.FriendListResp{List: friends}, nil
}

package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"im/apps/social/rpc/internal/config"
	"im/apps/social/socialmodels"
)

type ServiceContext struct {
	Config config.Config

	socialmodels.FriendsModel
	socialmodels.FriendRequestsModel
	socialmodels.GroupsModel
	socialmodels.GroupRequestsModel
	socialmodels.GroupMembersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:              c,
		FriendsModel:        socialmodels.NewFriendsModel(conn, c.Cache),
		FriendRequestsModel: socialmodels.NewFriendRequestsModel(conn, c.Cache),
		GroupsModel:         socialmodels.NewGroupsModel(conn, c.Cache),
		GroupRequestsModel:  socialmodels.NewGroupRequestsModel(conn, c.Cache),
		GroupMembersModel:   socialmodels.NewGroupMembersModel(conn, c.Cache),
	}
}

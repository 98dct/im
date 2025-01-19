package constants

type HandlerResult int

// 未处理：1  通过：2  拒绝：3  撤销：4
const (
	NoHandler HandlerResult = iota + 1
	Pass
	Reject
	Cancel
)

// 群等级 1. 创建者，2. 管理者，3. 普通
type GroupRoleLevel int

const (
	CreatorGroupRoleLevel GroupRoleLevel = iota + 1 // 为什么会 从1开始？
	ManagerGroupRoleLevel
	AtLargeGroupRoleLevel
)

// 进群申请的方式： 1. 邀请， 2. 申请
type GroupJoinSource int

const (
	InviteGroupJoinSource GroupJoinSource = iota + 1
	PutInGroupJoinSource
)

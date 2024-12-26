package constants

type HandlerResult int

// 未处理：1  通过：2  拒绝：3  撤销：4
const (
	NoHandler HandlerResult = iota + 1
	Pass
	Reject
	Cancel
)

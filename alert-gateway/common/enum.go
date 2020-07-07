package common

type AlertStatusEnum int

const (
	// 已被处理告警状态
	AlertStatusEnumOff AlertStatusEnum = 0
	// 已被确认告警状态
	AlertStatusEnumConfirm AlertStatusEnum = 1
	// 正在告警中告警状态
	AlertStatusEnumOn AlertStatusEnum = 2
)

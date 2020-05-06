package tencentcloud

const monitorEventTypeStatusChange = "status_change"
const monitorEventTypeAbnormal = "abnormal"

var monitorEventTypes = []string{
	monitorEventTypeStatusChange,
	monitorEventTypeAbnormal,
}

const monitorEventStatusRecover = "recover"
const monitorEventStatusAlarm = "alarm"
const monitorEventStatusNothing = "-"

var monitorEventStatus = []string{
	monitorEventStatusRecover,
	monitorEventStatusAlarm,
	monitorEventStatusNothing,
}

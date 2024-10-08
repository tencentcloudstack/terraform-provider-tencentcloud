package cdwdoris

const (
	INSTANCE_STATE_INIT      = "Init"
	INSTANCE_STATE_SERVING   = "Serving"
	INSTANCE_STATE_DELETED   = "Deleted"
	INSTANCE_STATE_ISOLATED  = "Isolated"
	INSTANCE_STATE_CHANGING  = "Changing"
	INSTANCE_STATE_UPGRADING = "Upgrading"
)

const (
	CHARGE_TYPE_PREPAID          = "PREPAID"
	CHARGE_TYPE_POSTPAID_BY_HOUR = "POSTPAID_BY_HOUR"
)

var CHARGE_TYPE = []string{
	CHARGE_TYPE_PREPAID,
	CHARGE_TYPE_POSTPAID_BY_HOUR,
}

const (
	WORKLOAD_GROUP_STATUS_OPEN  = "open"
	WORKLOAD_GROUP_STATUS_CLOSE = "close"
)

var WORKLOAD_GROUP_STATUS = []string{
	WORKLOAD_GROUP_STATUS_OPEN,
	WORKLOAD_GROUP_STATUS_CLOSE,
}

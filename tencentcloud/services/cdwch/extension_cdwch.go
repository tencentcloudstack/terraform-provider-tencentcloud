package cdwch

import cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"

const (
	CDWCH_PAY_MODE_HOUR   = "hour"
	CDWCH_PAY_MODE_PREPAY = "prepay"
)

const (
	CDWCH_CHARGE_TYPE_PREPAID          = "PREPAID"
	CDWCH_CHARGE_TYPE_POSTPAID_BY_HOUR = "POSTPAID_BY_HOUR"
)

const (
	NODE_TYPE_CLICKHOUSE = "DATA"
	NODE_TYPE_ZOOKEEPER  = "COMMON"
)

const (
	OPERATION_TYPE_CREATE = "create"
	OPERATION_TYPE_UPDATE = "update"
)

const (
	SCHEDULE_TYPE_DATA = "data"
	SCHEDULE_TYPE_META = "meta"
)

const (
	ACTION_ALTER_CK_USER_ADD_SYSTEM_USER    = "AddSystemUser"
	ACTION_ALTER_CK_USER_UPDATE_SYSTEM_USER = "UpdateSystemUser"
)

const (
	DESCRIBE_CK_SQL_APIS_GET_SYSTEM_USERS    = "GetSystemUsers"
	DESCRIBE_CK_SQL_APIS_REVOKE_CLUSTER_USER = "RevokeClusterUser"
	DESCRIBE_CK_SQL_APIS_DELETE_SYSTEM_USER  = "DeleteSystemUser"
)

var PAY_MODE_TO_CHARGE_TYPE = map[string]string{
	CDWCH_PAY_MODE_HOUR:   CDWCH_CHARGE_TYPE_POSTPAID_BY_HOUR,
	CDWCH_PAY_MODE_PREPAY: CDWCH_CHARGE_TYPE_PREPAID,
}

type AccountInfo struct {
	InstanceId string `json:"InstanceId"`
	UserName   string `json:"UserName"`
	Describe   string `json:"Describe"`
	Type       string `json:"Type"`
	Cluster    string `json:"Cluster"`
}

type AccountPermission struct {
	InstanceId            string                         `json:"InstanceId"`
	Cluster               string                         `json:"Cluster"`
	UserName              string                         `json:"UserName"`
	AllDatabase           bool                           `json:"AllDatabase"`
	GlobalPrivileges      []string                       `json:"GlobalPrivileges"`
	DatabasePrivilegeList []*cdwch.DatabasePrivilegeInfo `json:"DatabasePrivilegeList"`
}

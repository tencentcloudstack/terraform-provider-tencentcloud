package tencentcloud

const (
	CKAFKA_DESCRIBE_LIMIT    = 50
	CKAFKA_ACL_PRINCIPAL_STR = "User:"
)

var CKAFKA_ACL_RESOURCE_TYPE = map[string]int64{
	"UNKNOWN":          0,
	"ANY":              1,
	"TOPIC":            2,
	"GROUP":            3,
	"CLUSTER":          4,
	"TRANSACTIONAL_ID": 5,
}

var CKAFKA_ACL_RESOURCE_TYPE_TO_STRING = map[int64]string{
	0: "UNKNOWN",
	1: "ANY",
	2: "TOPIC",
	3: "GROUP",
	4: "CLUSTER",
	5: "TRANSACTIONAL_ID",
}

var CKAFKA_ACL_OPERATION = map[string]int64{
	"UNKNOWN":          0,
	"ANY":              1,
	"ALL":              2,
	"READ":             3,
	"WRITE":            4,
	"CREATE":           5,
	"DELETE":           6,
	"ALTER":            7,
	"DESCRIBE":         8,
	"CLUSTER_ACTION":   9,
	"DESCRIBE_CONFIGS": 10,
	"ALTER_CONFIGS":    11,
	"IDEMPOTEN_WRITE":  12,
}
var CKAFKA_ACL_OPERATION_TO_STRING = map[int64]string{
	0:  "UNKNOWN",
	1:  "ANY",
	2:  "ALL",
	3:  "READ",
	4:  "WRITE",
	5:  "CREATE",
	6:  "DELETE",
	7:  "ALTER",
	8:  "DESCRIBE",
	9:  "CLUSTER_ACTION",
	10: "DESCRIBE_CONFIGS",
	11: "ALTER_CONFIGS",
	12: "IDEMPOTEN_WRITE",
}

var CKAFKA_PERMISSION_TYPE = map[string]int64{
	"UNKNOWN": 0,
	"ANY":     1,
	"DENY":    2,
	"ALLOW":   3,
}

var CKAFKA_PERMISSION_TYPE_TO_STRING = map[int64]string{
	0: "UNKNOWN",
	1: "ANY",
	2: "DENY",
	3: "ALLOW",
}

//sdk ckafka not found error
const CkafkaInstanceNotFound = "InvalidParameterValue.InstanceNotExist"
const CkafkaFailedOperation = "FailedOperation"

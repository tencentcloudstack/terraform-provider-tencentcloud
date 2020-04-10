package tencentcloud

//redis version  https://cloud.tencent.com/document/api/239/20022#ProductConf
const (
	REDIS_VERSION_MASTER_SLAVE_REDIS  = 2
	REDIS_VERSION_MASTER_SLAVE_CKV    = 3
	REDIS_VERSION_CLUSTER_CKV         = 4
	REDIS_VERSION_STANDALONE_REDIS    = 5
	REDIS_VERSION_CLUSTER_REDIS       = 7
	REDIS_VERSION_MASTER_SLAVE_REDIS5 = 8
	REDIS_VERSION_CLUSTER_REDIS5      = 9
)

var REDIS_NAMES = map[int64]string{
	REDIS_VERSION_MASTER_SLAVE_REDIS:  "master_slave_redis",
	REDIS_VERSION_MASTER_SLAVE_CKV:    "master_slave_ckv",
	REDIS_VERSION_CLUSTER_REDIS:       "cluster_redis",
	REDIS_VERSION_CLUSTER_CKV:         "cluster_ckv",
	REDIS_VERSION_STANDALONE_REDIS:    "standalone_redis",
	REDIS_VERSION_MASTER_SLAVE_REDIS5: "master_slave_redis5.0",
	REDIS_VERSION_CLUSTER_REDIS5:      "cluster_redis5.0",
}

//redis status  https://cloud.tencent.com/document/product/239/20018
const (
	REDIS_STATUS_INIT       = 0
	REDIS_STATUS_PROCESSING = 1
	REDIS_STATUS_ONLINE     = 2
	REDIS_STATUS_ISOLATE    = -2
	REDIS_STATUS_TODELETE   = -3
)

var REDIS_STATUS = map[int64]string{
	REDIS_STATUS_INIT:       "init",
	REDIS_STATUS_PROCESSING: "processing",
	REDIS_STATUS_ONLINE:     "online",
	REDIS_STATUS_ISOLATE:    "isolate",
	REDIS_STATUS_TODELETE:   "todelete",
}

/*
	https://cloud.tencent.com/document/api/239/20022#TradeDealDetail
	Order status
	1: unpaid
	2: paid, not shipped
	3: in shipment
	4: successfully
	5: shipped failed
	6: refunded
	7: closed order
	8: expired
	9: order no longer valid
	10: product no longer valid
	11: payment refused
	12: in payment
*/
const (
	REDIS_ORDER_SUCCESS_DELIVERY = 4
	REDIS_ORDER_PAYMENT          = 12
)

//https://cloud.tencent.com/document/api/239/30601
const (
	REDIS_TASK_PREPARING = "preparing"
	REDIS_TASK_RUNNING   = "running"
	REDIS_TASK_SUCCEED   = "succeed"
	REDIS_TASK_FAILED    = "failed"
	REDIS_TASK_ERROR     = "error"
)

//sdk redis not found error
const RedisInstanceNotFound = "ResourceNotFound.InstanceNotExists"

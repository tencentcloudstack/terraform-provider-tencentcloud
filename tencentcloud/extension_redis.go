package tencentcloud

//redis version  https://cloud.tencent.com/document/api/239/20022#ProductConf
const (
	REDIS_VERSION_MASTER_SLAVE_REDIS = 2
	REDIS_VERSION_MASTER_SLAVE_CKV   = 3
	REDIS_VERSION_CLUSTER_CKV        = 4
	REDIS_VERSION_STANDALONE_REDIS   = 5
	REDIS_VERSION_CLUSTER_REDIS      = 7
)

var REDIS_NAMES = map[int64]string{
	REDIS_VERSION_MASTER_SLAVE_REDIS: "master_slave_redis",
	REDIS_VERSION_MASTER_SLAVE_CKV:   "master_slave_ckv",
	REDIS_VERSION_CLUSTER_REDIS:      "cluster_redis",
	REDIS_VERSION_CLUSTER_CKV:        "cluster_ckv",
	REDIS_VERSION_STANDALONE_REDIS:   "standalone_redis",
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

var REDIS_ZONE_ID2NAME = map[int64]string{
	100001: "ap-guangzhou-1",
	100002: "ap-guangzhou-2",
	100003: "ap-guangzhou-3",
	100004: "ap-guangzhou-4",
	200001: "ap-shanghai-1",
	200002: "ap-shanghai-2",
	200003: "ap-shanghai-3",
	200004: "ap-shanghai-4",
	300001: "ap-hongkong-1",
	300002: "ap-hongkong-2",
	400001: "na-toronto-1",
	700001: "ap-shanghai-fsi-1",
	700002: "ap-shanghai-fsi-2",
	800001: "ap-beijing-1",
	800002: "ap-beijing-2",
	800003: "ap-beijing-3",
	800004: "ap-beijing-4",
	900001: "ap-singapore-1",
	110001: "ap-shenzhen-fsi-1",
	110002: "ap-shenzhen-fsi-2",
	150001: "na-siliconvalley-1",
	150002: "na-siliconvalley-2",
	160001: "ap-chengdu-1",
	160002: "ap-chengdu-2",
	170001: "en-frankfurt-1",
	180001: "ap-seoul-1",
	190001: "ap-chongqing-1",
	210001: "ap-mumbai-1",
	220001: "na-ashburn-1",
	230001: "ap-bangkok-1",
	240001: "eu-moscow-1",
	250001: "ap-tokyo-1",
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

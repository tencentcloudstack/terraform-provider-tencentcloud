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

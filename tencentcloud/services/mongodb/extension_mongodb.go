package mongodb

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
)

const (
	MONGODB_INSTANCE_STATUS_INITIAL    = 0
	MONGODB_INSTANCE_STATUS_PROCESSING = 1
	MONGODB_INSTANCE_STATUS_RUNNING    = 2
	MONGODB_INSTANCE_STATUS_EXPIRED    = -2

	MONGODB_ENGINE_VERSION_3_WT    = "MONGO_3_WT"
	MONGODB_ENGINE_VERSION_36_WT   = "MONGO_36_WT"
	MONGODB_ENGINE_VERSION_3_ROCKS = "MONGO_3_ROCKS"
	MONGODB_ENGINE_VERSION_4_WT    = "MONGO_40_WT"

	MONGODB_MACHINE_TYPE_GIO    = "GIO"
	MONGODB_MACHINE_TYPE_TGIO   = "TGIO"
	MONGODB_MACHINE_TYPE_HIO    = "HIO"
	MONGODB_MACHINE_TYPE_HIO10G = "HIO10G"

	MONGODB_CLUSTER_TYPE_REPLSET = "REPLSET"
	MONGODB_CLUSTER_TYPE_SHARD   = "SHARD"

	MONGO_INSTANCE_TYPE_FORMAL   = 1
	MONGO_INSTANCE_TYPE_READONLY = 3
	MONGO_INSTANCE_TYPE_STANDBY  = 4
)

var MONGODB_CLUSTER_TYPE = []string{
	MONGODB_CLUSTER_TYPE_REPLSET,
	MONGODB_CLUSTER_TYPE_SHARD,
}

const (
	MONGODB_DEFAULT_LIMIT  = 20
	MONGODB_MAX_LIMIT      = 100
	MONGODB_DEFAULT_OFFSET = 0
)

const (
	MONGODB_CHARGE_TYPE_POSTPAID = svcpostgresql.COMMON_PAYTYPE_POSTPAID
	MONGODB_CHARGE_TYPE_PREPAID  = svcpostgresql.COMMON_PAYTYPE_PREPAID
)

var MONGODB_CHARGE_TYPE = map[uint64]string{
	0: MONGODB_CHARGE_TYPE_POSTPAID,
	1: MONGODB_CHARGE_TYPE_PREPAID,
}

var MONGODB_AUTO_RENEW_FLAG = map[int]string{
	0: "NOTIFY_AND_MANUAL_RENEW",
	1: "NOTIFY_AND_AUTO_RENEW",
	2: "DISABLE_NOTIFY_AND_MANUAL_RENEW",
}

var MONGODB_PREPAID_PERIOD = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36}

const (
	MONGODB_TASK_FAILED  = "failed"
	MONGODB_TASK_PAUSED  = "paused"
	MONGODB_TASK_RUNNING = "running"
	MONGODB_TASK_SUCCESS = "success"
)

const (
	MONGODB_STATUS_DELIVERY_SUCCESS = 4
)

func TencentMongodbBasicInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: tccommon.ValidateStringLengthInRange(2, 35),
			Description:  "Name of the Mongodb instance.",
		},
		"memory": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: tccommon.ValidateIntegerMin(2),
			Description:  "Memory size. The minimum value is 2, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.",
		},
		"volume": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: tccommon.ValidateIntegerMin(25),
			Description:  "Disk size. The minimum value is 25, and unit is GB. Memory and volume must be upgraded or degraded simultaneously.",
		},
		"engine_version": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Version of the Mongodb, and available values include `MONGO_36_WT` (MongoDB 3.6 WiredTiger Edition), `MONGO_40_WT` (MongoDB 4.0 WiredTiger Edition) and `MONGO_42_WT`  (MongoDB 4.2 WiredTiger Edition). NOTE: `MONGO_3_WT` (MongoDB 3.2 WiredTiger Edition) and `MONGO_3_ROCKS` (MongoDB 3.2 RocksDB Edition) will deprecated.",
		},
		"machine_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
				if (olds == MONGODB_MACHINE_TYPE_GIO && news == MONGODB_MACHINE_TYPE_HIO) ||
					(olds == MONGODB_MACHINE_TYPE_HIO && news == MONGODB_MACHINE_TYPE_GIO) {
					return true
				} else if (olds == MONGODB_MACHINE_TYPE_TGIO && news == MONGODB_MACHINE_TYPE_HIO10G) ||
					(olds == MONGODB_MACHINE_TYPE_HIO10G && news == MONGODB_MACHINE_TYPE_TGIO) {
					return true
				}
				return olds == news
			},
			Description: "Type of Mongodb instance, and available values include `HIO`(or `GIO` which will be deprecated, represents high IO) and `HIO10G`(or `TGIO` which will be deprecated, represents 10-gigabit high IO).",
		},
		"available_zone": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The available zone of the Mongodb.",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "",
			Description: "ID of the VPC.",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "ID of the subnet within this VPC. The value is required if `vpc_id` is set.",
		},
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "ID of the project which the instance belongs.",
		},
		"security_groups": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Set: func(v interface{}) int {
				return helper.HashString(v.(string))
			},
			Description: "ID of the security group.",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Password of this Mongodb account.",
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The tags of the Mongodb. Key name `project` is system reserved and can't be used.",
		},
		"mongos_cpu": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Number of mongos cpu.",
		},
		"mongos_memory": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Mongos memory size in GB.",
		},
		"mongos_node_num": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Number of mongos.",
		},
		// payment
		"charge_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      MONGODB_CHARGE_TYPE_POSTPAID,
			ValidateFunc: tccommon.ValidateAllowedStringValue([]string{MONGODB_CHARGE_TYPE_POSTPAID, MONGODB_CHARGE_TYPE_PREPAID}),
			Description:  "The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`. Caution that update operation on this field will delete old instances and create new one with new charge type.",
		},
		"prepaid_period": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: tccommon.ValidateAllowedIntValue(MONGODB_PREPAID_PERIOD),
			Description:  "The tenancy (time unit is month) of the prepaid instance. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. NOTE: it only works when charge_type is set to `PREPAID`.",
		},
		"auto_renew_flag": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Auto renew flag. Valid values are `0`(NOTIFY_AND_MANUAL_RENEW), `1`(NOTIFY_AND_AUTO_RENEW) and `2`(DISABLE_NOTIFY_AND_MANUAL_RENEW). Default value is `0`. Note: only works for PREPAID instance. Only supports`0` and `1` for creation.",
		},
		"in_maintenance": {
			Type:     schema.TypeInt,
			Optional: true,
			Description: "Switch time for instance configuration changes.\n" +
				"	- 0: When the adjustment is completed, perform the configuration task immediately. Default is 0.\n" +
				"	- 1: Perform reconfiguration tasks within the maintenance time window.\n" +
				"Note: Adjusting the number of nodes and slices does not support changes within the maintenance window.",
		},
		// Computed
		"status": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Status of the Mongodb instance, and available values include pending initialization(expressed with 0),  processing(expressed with 1), running(expressed with 2) and expired(expressed with -2).",
		},
		"vip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP of the Mongodb instance.",
		},
		"vport": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "IP port of the Mongodb instance.",
		},
		"create_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Creation time of the Mongodb instance.",
		},
	}
}

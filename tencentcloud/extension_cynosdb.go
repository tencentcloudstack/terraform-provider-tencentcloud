package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	CYNOSDB_CHARGE_TYPE_POSTPAID = COMMON_PAYTYPE_POSTPAID
	CYNOSDB_CHARGE_TYPE_PREPAID  = COMMON_PAYTYPE_PREPAID

	CYNOSDB_STATUS_RUNNING  = "running"
	CYNOSDB_STATUS_OFFLINE  = "offlined"
	CYNOSDB_STATUS_ISOLATED = "isolated"
	CYNOSDB_STATUS_DELETED  = "deleted"

	CYNOSDB_UPGRADE_IMMEDIATE = "upgradeImmediate"

	CYNOSDB_INSTANCE_RW_TYPE = "rw"
	CYNOSDB_INSTANCE_RO_TYPE = "ro"

	CYNOSDB_DEFAULT_OFFSET = 0
	CYNOSDB_MAX_LIMIT      = 100

	CYNOSDB_INSGRP_HA = "ha"
	CYNOSDB_INSGRP_RO = "ro"
)

var (
	CYNOSDB_PREPAID_PERIOD = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36}
	CYNOSDB_CHARGE_TYPE    = map[int64]string{
		0: COMMON_PAYTYPE_POSTPAID,
		1: COMMON_PAYTYPE_PREPAID,
	}
)

func TencentCynosdbInstanceBaseInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_cpu_core": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The number of CPU cores of read-write type instance in the CynosDB cluster. Note: modification of this field will take effect immediately, if want to upgrade on maintenance window, please upgrade from console.",
		},
		"instance_memory_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Memory capacity of read-write type instance, unit in GB. Note: modification of this field will take effect immediately, if want to upgrade on maintenance window, please upgrade from console.",
		},
		"instance_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of instance.",
		},
		"instance_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of instance.",
		},
		"instance_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Status of the instance.",
		},
		"instance_storage_size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Storage size of the instance, unit in GB.",
		},
		"instance_maintain_weekdays": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			// DefaultFunc doesn't work but wil remain it
			DefaultFunc: func() (interface{}, error) {
				weekdays := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
				return weekdays, nil
			},
			Elem: &schema.Schema{Type: schema.TypeString},
			Set: func(v interface{}) int {
				return hashcode.String(v.(string))
			},
			Description: "Weekdays for maintenance. `[\"Mon\", \"Tue\", \"Wed\", \"Thu\", \"Fri\", \"Sat\", \"Sun\"]` by default.",
		},
		"instance_maintain_start_time": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     10800,
			Description: "Offset time from 00:00, unit in second. For example, 03:00am should be `10800`. `10800` by default.",
		},
		"instance_maintain_duration": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     3600,
			Description: "Duration time for maintenance, unit in second. `3600` by default.",
		},
	}
}

func TencentCynosdbClusterBaseInfo() map[string]*schema.Schema {
	cluster := map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     0,
			Description: "ID of the project. `0` by default.",
		},
		"available_zone": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The available zone of the CynosDB Cluster.",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the VPC.",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the subnet within this VPC.",
		},
		"port": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     5432,
			Description: "Port of CynosDB cluster.",
		},
		"db_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Type of CynosDB, and available values include `MYSQL`.",
		},
		"db_version": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Version of CynosDB, which is related to `db_type`. For `MYSQL`, available value is `5.7`.",
		},
		"storage_limit": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "Storage limit of CynosDB cluster instance, unit in GB.",
		},
		"cluster_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of CynosDB cluster.",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Sensitive:   true,
			Description: "Password of `root` account.",
		},
		// payment
		"charge_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      CYNOSDB_CHARGE_TYPE_POSTPAID,
			ValidateFunc: validateAllowedStringValue([]string{CYNOSDB_CHARGE_TYPE_POSTPAID, CYNOSDB_CHARGE_TYPE_PREPAID}),
			Description:  "The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`.",
		},
		"prepaid_period": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateAllowedIntValue(CYNOSDB_PREPAID_PERIOD),
			Description:  "The tenancy (time unit is month) of the prepaid instance. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`. NOTE: it only works when charge_type is set to `PREPAID`.",
		},
		"auto_renew_flag": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     0,
			Description: "Auto renew flag. Valid values are `0`(MANUAL_RENEW), `1`(AUTO_RENEW). Default value is `0`. Only works for PREPAID cluster.",
		},
		"force_delete": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Indicate whether to delete cluster instance directly or not. Default is false. If set true, the cluster and its `All RELATED INSTANCES` will be deleted instead of staying recycle bin. Note: works for both `PREPAID` and `POSTPAID_BY_HOUR` cluster.",
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The tags of the CynosDB cluster.",
		},
		// Computed
		"charset": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Charset used by CynosDB cluster.",
		},
		"cluster_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Status of the Cynosdb cluster.",
		},
		"create_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Creation time of the CynosDB cluster.",
		},
		"storage_used": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Used storage of CynosDB cluster, unit in MB.",
		},
		// rw instance group infos
		"rw_group_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of read-write instance group.",
		},
		"rw_group_instances": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of instances in the read-write instance group.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"instance_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "ID of instance.",
					},
					"instance_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Name of instance.",
					},
				},
			},
		},
		"rw_group_sg": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "IDs of security group for `rw_group`.",
		},
		"rw_group_addr": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Read-write addresses. Each element contains the following attributes:",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ip": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "IP address for read-write connection.",
					},
					"port": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Port number for read-write connection.",
					},
				},
			},
		},
		// ro instance group infos
		"ro_group_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of read-only instance group.",
		},
		"ro_group_instances": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of instances in the read-only instance group.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"instance_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "ID of instance.",
					},
					"instance_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Name of instance.",
					},
				},
			},
		},
		"ro_group_sg": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "IDs of security group for `ro_group`.",
		},
		"ro_group_addr": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Readonly addresses. Each element contains the following attributes:",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ip": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "IP address for readonly connection.",
					},
					"port": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Port number for readonly connection.",
					},
				},
			},
		},
	}

	for k, v := range TencentCynosdbInstanceBaseInfo() {
		cluster[k] = v
	}

	return cluster
}

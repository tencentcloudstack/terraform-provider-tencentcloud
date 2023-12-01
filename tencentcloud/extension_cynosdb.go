package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

const (
	CYNOSDB_CHARGE_TYPE_POSTPAID = COMMON_PAYTYPE_POSTPAID
	CYNOSDB_CHARGE_TYPE_PREPAID  = COMMON_PAYTYPE_PREPAID
	CYNOSDB_SERVERLESS           = "SERVERLESS"

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

	// 0-成功，1-失败，2-处理中
	CYNOSDB_FLOW_STATUS_SUCCESSFUL = "0"
	CYNOSDB_FLOW_STATUS_FAILED     = "1"
	CYNOSDB_FLOW_STATUS_PROCESSING = "2"
)

const (
	STATUS_YES = "yes"
	STATUS_NO  = "no"

	RW_TYPE = "READWRITE"
	RO_TYPE = "READONLY"
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
			Optional:    true,
			Description: "The number of CPU cores of read-write type instance in the CynosDB cluster. Required while creating normal cluster. Note: modification of this field will take effect immediately, if want to upgrade on maintenance window, please upgrade from console.",
		},
		"instance_memory_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Memory capacity of read-write type instance, unit in GB. Required while creating normal cluster. Note: modification of this field will take effect immediately, if want to upgrade on maintenance window, please upgrade from console.",
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
				return helper.HashString(v.(string))
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
			Description: "ID of the VPC.",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "ID of the subnet within this VPC.",
		},
		"old_ip_reserve_hours": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Recycling time of the old address, must be filled in when modifying the vpcRecycling time of the old address, must be filled in when modifying the vpc.",
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
			Optional:    true,
			Description: "Storage limit of CynosDB cluster instance, unit in GB. The maximum storage of a non-serverless instance in GB. NOTE: If db_type is `MYSQL` and charge_type is `PREPAID`, the value cannot exceed the maximum storage corresponding to the CPU and memory specifications, and the transaction mode is `order and pay`. when charge_type is `POSTPAID_BY_HOUR`, this argument is unnecessary.",
		},
		"storage_pay_mode": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Cluster storage billing mode, pay-as-you-go: `0`-yearly/monthly: `1`-The default is pay-as-you-go. When the DbType is MYSQL, when the cluster computing billing mode is post-paid (including DbMode is SERVERLESS), the storage billing mode can only be billing by volume; rollback and cloning do not support yearly subscriptions monthly storage.",
		},
		"cluster_name": {
			Type:        schema.TypeString,
			Required:    true,
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
		"param_items": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specify parameter list of database. It is valid when prarm_template_id is set in create cluster. Use `data.tencentcloud_mysql_default_params` to query available parameter details.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of param, e.g. `character_set_server`.",
					},
					"old_value": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Param old value, indicates the value which already set, this value is required when modifying current_value.",
					},
					"current_value": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Param expected value to set.",
					},
				},
			},
		},
		"prarm_template_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The ID of the parameter template.",
		},
		"db_mode": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify DB mode, only available when `db_type` is `MYSQL`. Values: `NORMAL` (Default), `SERVERLESS`.",
		},
		"min_cpu": {
			Optional:    true,
			Type:        schema.TypeFloat,
			Description: "Minimum CPU core count, required while `db_mode` is `SERVERLESS`, request DescribeServerlessInstanceSpecs for more reference.",
		},
		"max_cpu": {
			Optional:    true,
			Type:        schema.TypeFloat,
			Description: "Maximum CPU core count, required while `db_mode` is `SERVERLESS`, request DescribeServerlessInstanceSpecs for more reference.",
		},
		"auto_pause": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify whether the cluster can auto-pause while `db_mode` is `SERVERLESS`. Values: `yes` (default), `no`.",
		},
		"auto_pause_delay": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Specify auto-pause delay in second while `db_mode` is `SERVERLESS`. Value range: `[600, 691200]`. Default: `600`.",
		},
		"serverless_status_flag": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateAllowedStringValue([]string{"resume", "pause"}),
			Description:  "Specify whether to pause or resume serverless cluster. values: `resume`, `pause`.",
		},
		"serverless_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Serverless cluster status. NOTE: This is a readonly attribute, to modify, please set `serverless_status_flag`.",
		},
	}

	for k, v := range TencentCynosdbInstanceBaseInfo() {
		cluster[k] = v
	}

	return cluster
}

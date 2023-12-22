package mariadb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbDcnDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbDcnDetailRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"dcn_details": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "DCN synchronization details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region where the instance resides.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone where the instance resides.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance IP address.",
						},
						"vipv6": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance IPv6 address.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance port.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance status.",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance status description.",
						},
						"dcn_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DCN flag. Valid values: `1` (primary), `2` (disaster recovery).",
						},
						"dcn_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DCN status. Valid values: `0` (none), `1` (creating), `2` (syncing), `3` (disconnected).",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of CPU cores of the instance.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance memory capacity in GB.",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance storage capacity in GB.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing mode.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the instance in the format of 2006-01-02 15:04:05.",
						},
						"period_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expiration time of the instance in the format of 2006-01-02 15:04:05.",
						},
						"instance_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance type. Valid values: `1` (dedicated primary instance), `2` (non-dedicated primary instance), `3` (non-dedicated disaster recovery instance), `4` (dedicated disaster recovery instance).",
						},
						"replica_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Configuration information of DCN replication. This field is null for a primary instance.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ro_replication_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DCN running status. Valid values: `START` (running), `STOP` (pause)Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"delay_replication_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Delayed replication type. Valid values: `DEFAULT` (no delay), `DUE_TIME` (specified replication time)Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"due_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specified time for delayed replicationNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"replication_delay": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of seconds to delay the replicationNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"replica_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "DCN replication status. This field is null for the primary instance.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DCN running status. Valid values: `START` (running), `STOP` (pause).Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"delay": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The current delay, which takes the delay value of the replica instance.",
									},
								},
							},
						},
						"encrypt_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether KMS is enabled.",
						},
					},
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMariadbDcnDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_dcn_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dcnDetails []*mariadb.DcnDetailItem
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbDcnDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		dcnDetails = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(dcnDetails))
	tmpList := make([]map[string]interface{}, 0, len(dcnDetails))

	if dcnDetails != nil {
		for _, dcnDetailItem := range dcnDetails {
			dcnDetailItemMap := map[string]interface{}{}

			if dcnDetailItem.InstanceId != nil {
				dcnDetailItemMap["instance_id"] = dcnDetailItem.InstanceId
			}

			if dcnDetailItem.InstanceName != nil {
				dcnDetailItemMap["instance_name"] = dcnDetailItem.InstanceName
			}

			if dcnDetailItem.Region != nil {
				dcnDetailItemMap["region"] = dcnDetailItem.Region
			}

			if dcnDetailItem.Zone != nil {
				dcnDetailItemMap["zone"] = dcnDetailItem.Zone
			}

			if dcnDetailItem.Vip != nil {
				dcnDetailItemMap["vip"] = dcnDetailItem.Vip
			}

			if dcnDetailItem.Vipv6 != nil {
				dcnDetailItemMap["vipv6"] = dcnDetailItem.Vipv6
			}

			if dcnDetailItem.Vport != nil {
				dcnDetailItemMap["vport"] = dcnDetailItem.Vport
			}

			if dcnDetailItem.Status != nil {
				dcnDetailItemMap["status"] = dcnDetailItem.Status
			}

			if dcnDetailItem.StatusDesc != nil {
				dcnDetailItemMap["status_desc"] = dcnDetailItem.StatusDesc
			}

			if dcnDetailItem.DcnFlag != nil {
				dcnDetailItemMap["dcn_flag"] = dcnDetailItem.DcnFlag
			}

			if dcnDetailItem.DcnStatus != nil {
				dcnDetailItemMap["dcn_status"] = dcnDetailItem.DcnStatus
			}

			if dcnDetailItem.Cpu != nil {
				dcnDetailItemMap["cpu"] = dcnDetailItem.Cpu
			}

			if dcnDetailItem.Memory != nil {
				dcnDetailItemMap["memory"] = dcnDetailItem.Memory
			}

			if dcnDetailItem.Storage != nil {
				dcnDetailItemMap["storage"] = dcnDetailItem.Storage
			}

			if dcnDetailItem.PayMode != nil {
				dcnDetailItemMap["pay_mode"] = dcnDetailItem.PayMode
			}

			if dcnDetailItem.CreateTime != nil {
				dcnDetailItemMap["create_time"] = dcnDetailItem.CreateTime
			}

			if dcnDetailItem.PeriodEndTime != nil {
				dcnDetailItemMap["period_end_time"] = dcnDetailItem.PeriodEndTime
			}

			if dcnDetailItem.InstanceType != nil {
				dcnDetailItemMap["instance_type"] = dcnDetailItem.InstanceType
			}

			if dcnDetailItem.ReplicaConfig != nil {
				replicaConfigMap := map[string]interface{}{}

				if dcnDetailItem.ReplicaConfig.RoReplicationMode != nil {
					replicaConfigMap["ro_replication_mode"] = dcnDetailItem.ReplicaConfig.RoReplicationMode
				}

				if dcnDetailItem.ReplicaConfig.DelayReplicationType != nil {
					replicaConfigMap["delay_replication_type"] = dcnDetailItem.ReplicaConfig.DelayReplicationType
				}

				if dcnDetailItem.ReplicaConfig.DueTime != nil {
					replicaConfigMap["due_time"] = dcnDetailItem.ReplicaConfig.DueTime
				}

				if dcnDetailItem.ReplicaConfig.ReplicationDelay != nil {
					replicaConfigMap["replication_delay"] = dcnDetailItem.ReplicaConfig.ReplicationDelay
				}

				dcnDetailItemMap["replica_config"] = []interface{}{replicaConfigMap}
			}

			if dcnDetailItem.ReplicaStatus != nil {
				replicaStatusMap := map[string]interface{}{}

				if dcnDetailItem.ReplicaStatus.Status != nil {
					replicaStatusMap["status"] = dcnDetailItem.ReplicaStatus.Status
				}

				if dcnDetailItem.ReplicaStatus.Delay != nil {
					replicaStatusMap["delay"] = dcnDetailItem.ReplicaStatus.Delay
				}

				dcnDetailItemMap["replica_status"] = []interface{}{replicaStatusMap}
			}

			if dcnDetailItem.EncryptStatus != nil {
				dcnDetailItemMap["encrypt_status"] = dcnDetailItem.EncryptStatus
			}

			ids = append(ids, *dcnDetailItem.InstanceId)
			tmpList = append(tmpList, dcnDetailItemMap)
		}

		_ = d.Set("dcn_details", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}

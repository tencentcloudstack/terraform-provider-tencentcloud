package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbClusterInstanceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClusterInstanceGroupsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of cluster.",
			},

			"instance_grp_info_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of instance groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "App id.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of cluster.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created time.",
						},
						"deleted_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deleted time.",
						},
						"instance_grp_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of instance group.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance group type. ha-ha group; ro-read-only group.",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updated time.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Intranet IP.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Intranet port.",
						},
						"wan_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public domain name.",
						},
						"wan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public IP.",
						},
						"wan_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Public port.",
						},
						"wan_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public status.",
						},
						"instance_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance groups contain instance information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User Uin.",
									},
									"app_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "User app id.",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The id of cluster.",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of cluster.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The id of instance.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of instance.",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The id of project.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region.",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability zone.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status of instance.",
									},
									"status_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance state Chinese description.",
									},
									"db_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database type.",
									},
									"db_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database version.",
									},
									"cpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cpu, unit: CORE.",
									},
									"memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Memory, unit: GB.",
									},
									"storage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Storage, unit: GB.",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance type.",
									},
									"instance_role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance role.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Update time.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Create time.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC network ID.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet ID.",
									},
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance intranet IP.",
									},
									"vport": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance intranet VPort.",
									},
									"pay_mode": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Pay mode.",
									},
									"period_end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance expiration time.",
									},
									"destroy_deadline_text": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Destroy deadline.",
									},
									"isolate_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Isolate time.",
									},
									"net_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Net type.",
									},
									"wan_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Public domain.",
									},
									"wan_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Public IP.",
									},
									"wan_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Public port.",
									},
									"wan_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Public status.",
									},
									"destroy_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance destroy time.",
									},
									"cynos_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cynos kernel version.",
									},
									"processing_task": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Task being processed.",
									},
									"renew_flag": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Renew flag.",
									},
									"min_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Serverless instance minimum cpu.",
									},
									"max_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Serverless instance maxmum cpu.",
									},
									"serverless_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Serverless instance status, optional values:resumepause.",
									},
									"storage_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Prepaid Storage Id.Note: This field may return null, indicating that no valid value can be obtained..",
									},
									"storage_pay_mode": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Storage payment type.",
									},
									"physical_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Physical zone.",
									},
									"business_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Business type.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"tasks": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Task list.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"task_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Task auto-increment ID.Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"task_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Task type.Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"task_status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Task status.Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"object_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Task ID (cluster ID|instance group ID|instance ID).Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"object_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Object type.Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"is_freeze": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to freeze.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"resource_tags": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Resource tags.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tag_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The key of tag.",
												},
												"tag_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The value of tag.",
												},
											},
										},
									},
								},
							},
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

func dataSourceTencentCloudCynosdbClusterInstanceGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_cluster_instance_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceGrpInfoList []*cynosdb.CynosdbInstanceGrp

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClusterInstanceGrps(ctx, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceGrpInfoList = result.Response.InstanceGrpInfoList
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceGrpInfoList))
	tmpList := make([]map[string]interface{}, 0, len(instanceGrpInfoList))
	for _, instanceGrpInfo := range instanceGrpInfoList {
		ids = append(ids, *instanceGrpInfo.InstanceGrpId)
		instanceGrpInfoMap := make(map[string]interface{})
		instanceGrpInfoMap["app_id"] = instanceGrpInfo.AppId
		instanceGrpInfoMap["cluster_id"] = instanceGrpInfo.ClusterId
		instanceGrpInfoMap["created_time"] = instanceGrpInfo.CreatedTime
		instanceGrpInfoMap["deleted_time"] = instanceGrpInfo.DeletedTime
		instanceGrpInfoMap["instance_grp_id"] = instanceGrpInfo.InstanceGrpId
		instanceGrpInfoMap["status"] = instanceGrpInfo.Status
		instanceGrpInfoMap["type"] = instanceGrpInfo.Type
		instanceGrpInfoMap["updated_time"] = instanceGrpInfo.UpdatedTime
		instanceGrpInfoMap["vip"] = instanceGrpInfo.Vip
		instanceGrpInfoMap["vport"] = instanceGrpInfo.Vport
		instanceGrpInfoMap["wan_domain"] = instanceGrpInfo.WanDomain
		instanceGrpInfoMap["wan_ip"] = instanceGrpInfo.WanIP
		instanceGrpInfoMap["wan_port"] = instanceGrpInfo.WanPort
		instanceGrpInfoMap["wan_status"] = instanceGrpInfo.WanStatus
		if instanceGrpInfo.InstanceSet != nil {
			instances := make([]map[string]interface{}, 0)
			for _, instance := range instanceGrpInfo.InstanceSet {
				instanceMap := make(map[string]interface{})
				instanceMap["uin"] = instance.Uin
				instanceMap["app_id"] = instance.AppId
				instanceMap["cluster_id"] = instance.ClusterId
				instanceMap["cluster_name"] = instance.ClusterName
				instanceMap["instance_id"] = instance.InstanceId
				instanceMap["instance_name"] = instance.InstanceName
				instanceMap["project_id"] = instance.ProjectId
				instanceMap["region"] = instance.Region
				instanceMap["zone"] = instance.Zone
				instanceMap["status"] = instance.Status
				instanceMap["status_desc"] = instance.StatusDesc
				instanceMap["db_type"] = instance.DbType
				instanceMap["db_version"] = instance.DbVersion
				instanceMap["cpu"] = instance.Cpu
				instanceMap["memory"] = instance.Memory
				instanceMap["storage"] = instance.Storage
				instanceMap["instance_type"] = instance.InstanceType
				instanceMap["instance_role"] = instance.InstanceRole
				instanceMap["update_time"] = instance.UpdateTime
				instanceMap["create_time"] = instance.CreateTime
				instanceMap["vpc_id"] = instance.VpcId
				instanceMap["subnet_id"] = instance.SubnetId
				instanceMap["vip"] = instance.Vip
				instanceMap["vport"] = instance.Vport
				instanceMap["pay_mode"] = instance.PayMode
				instanceMap["period_end_time"] = instance.PeriodEndTime
				instanceMap["destroy_deadline_text"] = instance.DestroyDeadlineText
				instanceMap["isolate_time"] = instance.IsolateTime
				instanceMap["net_type"] = instance.NetType
				instanceMap["wan_domain"] = instance.WanDomain
				instanceMap["wan_ip"] = instance.WanIP
				instanceMap["wan_port"] = instance.WanPort
				instanceMap["wan_status"] = instance.WanStatus
				instanceMap["destroy_time"] = instance.DestroyTime
				instanceMap["cynos_version"] = instance.CynosVersion
				instanceMap["processing_task"] = instance.ProcessingTask
				instanceMap["renew_flag"] = instance.RenewFlag
				instanceMap["min_cpu"] = instance.MinCpu
				instanceMap["max_cpu"] = instance.MaxCpu
				instanceMap["serverless_status"] = instance.ServerlessStatus
				instanceMap["storage_id"] = instance.StorageId
				instanceMap["storage_pay_mode"] = instance.StoragePayMode
				instanceMap["physical_zone"] = instance.PhysicalZone
				instanceMap["business_type"] = instance.BusinessType
				instanceMap["is_freeze"] = instance.IsFreeze
				tasks := make([]map[string]interface{}, 0)
				if instance.Tasks != nil {
					for _, task := range instance.Tasks {
						taskMap := make(map[string]interface{})
						taskMap["task_id"] = task.TaskId
						taskMap["task_type"] = task.TaskType
						taskMap["task_status"] = task.TaskStatus
						taskMap["object_id"] = task.ObjectId
						taskMap["object_type"] = task.ObjectType

						tasks = append(tasks, taskMap)
					}
					instanceMap["tasks"] = tasks
				}
				tags := make([]map[string]interface{}, 0)
				if instance.ResourceTags != nil {
					for _, tag := range instance.ResourceTags {
						tagMap := make(map[string]interface{})
						tagMap["tag_key"] = tag.TagKey
						tagMap["tag_value"] = tag.TagValue

						tags = append(tags, tagMap)
					}
					instanceMap["resource_tags"] = tags
				}
				instances = append(instances, instanceMap)
				instanceGrpInfoMap["instance_set"] = instances
			}
		}
		tmpList = append(tmpList, instanceGrpInfoMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("instance_grp_info_list", tmpList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}

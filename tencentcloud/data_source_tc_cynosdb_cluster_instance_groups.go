/*
Use this data source to query detailed information of cynosdb cluster_instance_groups

Example Usage

```hcl
data "tencentcloud_cynosdb_cluster_instance_groups" "cluster_instance_groups" {
  cluster_id = &lt;nil&gt;
  total_count = &lt;nil&gt;
  instance_grp_info_list {
		app_id = &lt;nil&gt;
		cluster_id = &lt;nil&gt;
		created_time = &lt;nil&gt;
		deleted_time = &lt;nil&gt;
		instance_grp_id = &lt;nil&gt;
		status = &lt;nil&gt;
		type = &lt;nil&gt;
		updated_time = &lt;nil&gt;
		vip = &lt;nil&gt;
		vport = &lt;nil&gt;
		wan_domain = &lt;nil&gt;
		wan_i_p = &lt;nil&gt;
		wan_port = &lt;nil&gt;
		wan_status = &lt;nil&gt;
		instance_set {
			uin = &lt;nil&gt;
			app_id = &lt;nil&gt;
			cluster_id = &lt;nil&gt;
			cluster_name = &lt;nil&gt;
			instance_id = &lt;nil&gt;
			instance_name = &lt;nil&gt;
			project_id = &lt;nil&gt;
			region = &lt;nil&gt;
			zone = &lt;nil&gt;
			status = &lt;nil&gt;
			status_desc = &lt;nil&gt;
			db_type = &lt;nil&gt;
			db_version = &lt;nil&gt;
			cpu = &lt;nil&gt;
			memory = &lt;nil&gt;
			storage = &lt;nil&gt;
			instance_type = &lt;nil&gt;
			instance_role = &lt;nil&gt;
			update_time = &lt;nil&gt;
			create_time = &lt;nil&gt;
			vpc_id = &lt;nil&gt;
			subnet_id = &lt;nil&gt;
			vip = &lt;nil&gt;
			vport = &lt;nil&gt;
			pay_mode = &lt;nil&gt;
			period_end_time = &lt;nil&gt;
			destroy_deadline_text = &lt;nil&gt;
			isolate_time = &lt;nil&gt;
			net_type = &lt;nil&gt;
			wan_domain = &lt;nil&gt;
			wan_i_p = &lt;nil&gt;
			wan_port = &lt;nil&gt;
			wan_status = &lt;nil&gt;
			destroy_time = &lt;nil&gt;
			cynos_version = &lt;nil&gt;
			processing_task = &lt;nil&gt;
			renew_flag = &lt;nil&gt;
			min_cpu = &lt;nil&gt;
			max_cpu = &lt;nil&gt;
			serverless_status = &lt;nil&gt;
			storage_id = &lt;nil&gt;
			storage_pay_mode = &lt;nil&gt;
			physical_zone = &lt;nil&gt;
			business_type = &lt;nil&gt;
			tasks {
				task_id = &lt;nil&gt;
				task_type = &lt;nil&gt;
				task_status = &lt;nil&gt;
				object_id = &lt;nil&gt;
				object_type = &lt;nil&gt;
			}
			is_freeze = &lt;nil&gt;
			resource_tags {
				tag_key = &lt;nil&gt;
				tag_value = &lt;nil&gt;
			}
		}

  }
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbClusterInstanceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClusterInstanceGroupsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of cluster.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Number of instance groups.",
			},

			"instance_grp_info_list": {
				Type:        schema.TypeList,
				Description: "List of instance groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeInt,
							Description: "App id.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Description: "The ID of cluster.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Description: "Created time.",
						},
						"deleted_time": {
							Type:        schema.TypeString,
							Description: "Deleted time.",
						},
						"instance_grp_id": {
							Type:        schema.TypeString,
							Description: "The ID of instance group.",
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status.",
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Instance group type. ha-ha group; ro-read-only group.",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Description: "Updated time.",
						},
						"vip": {
							Type:        schema.TypeString,
							Description: "Intranet IP.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Description: "Intranet port.",
						},
						"wan_domain": {
							Type:        schema.TypeString,
							Description: "Public domain name.",
						},
						"wan_i_p": {
							Type:        schema.TypeString,
							Description: "Public IP.",
						},
						"wan_port": {
							Type:        schema.TypeInt,
							Description: "Public port.",
						},
						"wan_status": {
							Type:        schema.TypeString,
							Description: "Public status.",
						},
						"instance_set": {
							Type:        schema.TypeList,
							Description: "Instance groups contain instance information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uin": {
										Type:        schema.TypeString,
										Description: "User Uin.",
									},
									"app_id": {
										Type:        schema.TypeInt,
										Description: "User app id.",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Description: "The id of cluster.",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Description: "The name of cluster.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Description: "The id of instance.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Description: "The name of instance.",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Description: "The id of project.",
									},
									"region": {
										Type:        schema.TypeString,
										Description: "Region.",
									},
									"zone": {
										Type:        schema.TypeString,
										Description: "Availability zone.",
									},
									"status": {
										Type:        schema.TypeString,
										Description: "The status of instance.",
									},
									"status_desc": {
										Type:        schema.TypeString,
										Description: "Instance state Chinese description.",
									},
									"db_type": {
										Type:        schema.TypeString,
										Description: "Database type.",
									},
									"db_version": {
										Type:        schema.TypeString,
										Description: "Database version.",
									},
									"cpu": {
										Type:        schema.TypeInt,
										Description: "Cpu，unit: CORE.",
									},
									"memory": {
										Type:        schema.TypeInt,
										Description: "Memory, unit: GB.",
									},
									"storage": {
										Type:        schema.TypeInt,
										Description: "Storage, unit: GB.",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Description: "Instance type.",
									},
									"instance_role": {
										Type:        schema.TypeString,
										Description: "Instance role.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Description: "Update time.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Description: "Create time.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Description: "VPC network ID.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Description: "Subnet ID.",
									},
									"vip": {
										Type:        schema.TypeString,
										Description: "Instance intranet IP.",
									},
									"vport": {
										Type:        schema.TypeInt,
										Description: "Instance intranet VPort.",
									},
									"pay_mode": {
										Type:        schema.TypeInt,
										Description: "Pay mode.",
									},
									"period_end_time": {
										Type:        schema.TypeString,
										Description: "Instance expiration time.",
									},
									"destroy_deadline_text": {
										Type:        schema.TypeString,
										Description: "Destroy deadline.",
									},
									"isolate_time": {
										Type:        schema.TypeString,
										Description: "Isolate time.",
									},
									"net_type": {
										Type:        schema.TypeInt,
										Description: "Net type.",
									},
									"wan_domain": {
										Type:        schema.TypeString,
										Description: "Public domain.",
									},
									"wan_i_p": {
										Type:        schema.TypeString,
										Description: "Public IP.",
									},
									"wan_port": {
										Type:        schema.TypeInt,
										Description: "Public port.",
									},
									"wan_status": {
										Type:        schema.TypeString,
										Description: "Public status.",
									},
									"destroy_time": {
										Type:        schema.TypeString,
										Description: "Instance destory time.",
									},
									"cynos_version": {
										Type:        schema.TypeString,
										Description: "Cynos kernel version.",
									},
									"processing_task": {
										Type:        schema.TypeString,
										Description: "Task being processed.",
									},
									"renew_flag": {
										Type:        schema.TypeInt,
										Description: "Renew flag.",
									},
									"min_cpu": {
										Type:        schema.TypeFloat,
										Description: "Serverless instance minimum cpu.",
									},
									"max_cpu": {
										Type:        schema.TypeFloat,
										Description: "Serverless instance maxmum cpu.",
									},
									"serverless_status": {
										Type:        schema.TypeString,
										Description: "Serverless instance status, optional values:resumepause.",
									},
									"storage_id": {
										Type:        schema.TypeString,
										Description: "Prepaid Storage Id.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"storage_pay_mode": {
										Type:        schema.TypeInt,
										Description: "Storage payment type.",
									},
									"physical_zone": {
										Type:        schema.TypeString,
										Description: "Physical zone.",
									},
									"business_type": {
										Type:        schema.TypeString,
										Description: "Business type.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"tasks": {
										Type:        schema.TypeList,
										Description: "Task list.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"task_id": {
													Type:        schema.TypeInt,
													Description: "Task auto-increment ID.Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"task_type": {
													Type:        schema.TypeString,
													Description: "Task type.Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"task_status": {
													Type:        schema.TypeString,
													Description: "Task status.Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"object_id": {
													Type:        schema.TypeString,
													Description: "Task ID (cluster ID|instance group ID|instance ID).Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"object_type": {
													Type:        schema.TypeString,
													Description: "Object type.Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"is_freeze": {
										Type:        schema.TypeString,
										Description: "Whether to freeze.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"resource_tags": {
										Type:        schema.TypeList,
										Description: "Resource tags.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tag_key": {
													Type:        schema.TypeString,
													Description: "The key of tag.",
												},
												"tag_value": {
													Type:        schema.TypeString,
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
	defer logElapsed("data_source.tencentcloud_cynosdb_cluster_instance_groups.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_grp_info_list"); ok {
		instanceGrpInfoListSet := v.([]interface{})
		tmpSet := make([]*cynosdb.CynosdbInstanceGrp, 0, len(instanceGrpInfoListSet))

		for _, item := range instanceGrpInfoListSet {
			cynosdbInstanceGrp := cynosdb.CynosdbInstanceGrp{}
			cynosdbInstanceGrpMap := item.(map[string]interface{})

			if v, ok := cynosdbInstanceGrpMap["app_id"]; ok {
				cynosdbInstanceGrp.AppId = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbInstanceGrpMap["cluster_id"]; ok {
				cynosdbInstanceGrp.ClusterId = helper.String(v.(string))
			}
			if v, ok := cynosdbInstanceGrpMap["created_time"]; ok {
				cynosdbInstanceGrp.CreatedTime = helper.String(v.(string))
			}
			if v, ok := cynosdbInstanceGrpMap["deleted_time"]; ok {
				cynosdbInstanceGrp.DeletedTime = helper.String(v.(string))
			}
			if v, ok := cynosdbInstanceGrpMap["instance_grp_id"]; ok {
				cynosdbInstanceGrp.InstanceGrpId = helper.String(v.(string))
			}
			if v, ok := cynosdbInstanceGrpMap["status"]; ok {
				cynosdbInstanceGrp.Status = helper.String(v.(string))
			}
			if v, ok := cynosdbInstanceGrpMap["type"]; ok {
				cynosdbInstanceGrp.Type = helper.String(v.(string))
			}
			if v, ok := cynosdbInstanceGrpMap["updated_time"]; ok {
				cynosdbInstanceGrp.UpdatedTime = helper.String(v.(string))
			}
			if v, ok := cynosdbInstanceGrpMap["vip"]; ok {
				cynosdbInstanceGrp.Vip = helper.String(v.(string))
			}
			if v, ok := cynosdbInstanceGrpMap["vport"]; ok {
				cynosdbInstanceGrp.Vport = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbInstanceGrpMap["wan_domain"]; ok {
				cynosdbInstanceGrp.WanDomain = helper.String(v.(string))
			}
			if v, ok := cynosdbInstanceGrpMap["wan_i_p"]; ok {
				cynosdbInstanceGrp.WanIP = helper.String(v.(string))
			}
			if v, ok := cynosdbInstanceGrpMap["wan_port"]; ok {
				cynosdbInstanceGrp.WanPort = helper.IntInt64(v.(int))
			}
			if v, ok := cynosdbInstanceGrpMap["wan_status"]; ok {
				cynosdbInstanceGrp.WanStatus = helper.String(v.(string))
			}
			if v, ok := cynosdbInstanceGrpMap["instance_set"]; ok {
				for _, item := range v.([]interface{}) {
					instanceSetMap := item.(map[string]interface{})
					cynosdbInstance := cynosdb.CynosdbInstance{}
					if v, ok := instanceSetMap["uin"]; ok {
						cynosdbInstance.Uin = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["app_id"]; ok {
						cynosdbInstance.AppId = helper.IntInt64(v.(int))
					}
					if v, ok := instanceSetMap["cluster_id"]; ok {
						cynosdbInstance.ClusterId = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["cluster_name"]; ok {
						cynosdbInstance.ClusterName = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["instance_id"]; ok {
						cynosdbInstance.InstanceId = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["instance_name"]; ok {
						cynosdbInstance.InstanceName = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["project_id"]; ok {
						cynosdbInstance.ProjectId = helper.IntInt64(v.(int))
					}
					if v, ok := instanceSetMap["region"]; ok {
						cynosdbInstance.Region = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["zone"]; ok {
						cynosdbInstance.Zone = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["status"]; ok {
						cynosdbInstance.Status = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["status_desc"]; ok {
						cynosdbInstance.StatusDesc = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["db_type"]; ok {
						cynosdbInstance.DbType = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["db_version"]; ok {
						cynosdbInstance.DbVersion = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["cpu"]; ok {
						cynosdbInstance.Cpu = helper.IntInt64(v.(int))
					}
					if v, ok := instanceSetMap["memory"]; ok {
						cynosdbInstance.Memory = helper.IntInt64(v.(int))
					}
					if v, ok := instanceSetMap["storage"]; ok {
						cynosdbInstance.Storage = helper.IntInt64(v.(int))
					}
					if v, ok := instanceSetMap["instance_type"]; ok {
						cynosdbInstance.InstanceType = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["instance_role"]; ok {
						cynosdbInstance.InstanceRole = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["update_time"]; ok {
						cynosdbInstance.UpdateTime = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["create_time"]; ok {
						cynosdbInstance.CreateTime = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["vpc_id"]; ok {
						cynosdbInstance.VpcId = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["subnet_id"]; ok {
						cynosdbInstance.SubnetId = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["vip"]; ok {
						cynosdbInstance.Vip = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["vport"]; ok {
						cynosdbInstance.Vport = helper.IntInt64(v.(int))
					}
					if v, ok := instanceSetMap["pay_mode"]; ok {
						cynosdbInstance.PayMode = helper.IntInt64(v.(int))
					}
					if v, ok := instanceSetMap["period_end_time"]; ok {
						cynosdbInstance.PeriodEndTime = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["destroy_deadline_text"]; ok {
						cynosdbInstance.DestroyDeadlineText = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["isolate_time"]; ok {
						cynosdbInstance.IsolateTime = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["net_type"]; ok {
						cynosdbInstance.NetType = helper.IntInt64(v.(int))
					}
					if v, ok := instanceSetMap["wan_domain"]; ok {
						cynosdbInstance.WanDomain = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["wan_i_p"]; ok {
						cynosdbInstance.WanIP = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["wan_port"]; ok {
						cynosdbInstance.WanPort = helper.IntInt64(v.(int))
					}
					if v, ok := instanceSetMap["wan_status"]; ok {
						cynosdbInstance.WanStatus = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["destroy_time"]; ok {
						cynosdbInstance.DestroyTime = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["cynos_version"]; ok {
						cynosdbInstance.CynosVersion = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["processing_task"]; ok {
						cynosdbInstance.ProcessingTask = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["renew_flag"]; ok {
						cynosdbInstance.RenewFlag = helper.IntInt64(v.(int))
					}
					if v, ok := instanceSetMap["min_cpu"]; ok {
						cynosdbInstance.MinCpu = helper.Float64(v.(float64))
					}
					if v, ok := instanceSetMap["max_cpu"]; ok {
						cynosdbInstance.MaxCpu = helper.Float64(v.(float64))
					}
					if v, ok := instanceSetMap["serverless_status"]; ok {
						cynosdbInstance.ServerlessStatus = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["storage_id"]; ok {
						cynosdbInstance.StorageId = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["storage_pay_mode"]; ok {
						cynosdbInstance.StoragePayMode = helper.IntInt64(v.(int))
					}
					if v, ok := instanceSetMap["physical_zone"]; ok {
						cynosdbInstance.PhysicalZone = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["business_type"]; ok {
						cynosdbInstance.BusinessType = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["tasks"]; ok {
						for _, item := range v.([]interface{}) {
							tasksMap := item.(map[string]interface{})
							objectTask := cynosdb.ObjectTask{}
							if v, ok := tasksMap["task_id"]; ok {
								objectTask.TaskId = helper.IntInt64(v.(int))
							}
							if v, ok := tasksMap["task_type"]; ok {
								objectTask.TaskType = helper.String(v.(string))
							}
							if v, ok := tasksMap["task_status"]; ok {
								objectTask.TaskStatus = helper.String(v.(string))
							}
							if v, ok := tasksMap["object_id"]; ok {
								objectTask.ObjectId = helper.String(v.(string))
							}
							if v, ok := tasksMap["object_type"]; ok {
								objectTask.ObjectType = helper.String(v.(string))
							}
							cynosdbInstance.Tasks = append(cynosdbInstance.Tasks, &objectTask)
						}
					}
					if v, ok := instanceSetMap["is_freeze"]; ok {
						cynosdbInstance.IsFreeze = helper.String(v.(string))
					}
					if v, ok := instanceSetMap["resource_tags"]; ok {
						for _, item := range v.([]interface{}) {
							resourceTagsMap := item.(map[string]interface{})
							tag := cynosdb.Tag{}
							if v, ok := resourceTagsMap["tag_key"]; ok {
								tag.TagKey = helper.String(v.(string))
							}
							if v, ok := resourceTagsMap["tag_value"]; ok {
								tag.TagValue = helper.String(v.(string))
							}
							cynosdbInstance.ResourceTags = append(cynosdbInstance.ResourceTags, &tag)
						}
					}
					cynosdbInstanceGrp.InstanceSet = append(cynosdbInstanceGrp.InstanceSet, &cynosdbInstance)
				}
			}
			tmpSet = append(tmpSet, &cynosdbInstanceGrp)
		}
		paramMap["instance_grp_info_list"] = tmpSet
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceGrpInfoList []*cynosdb.CynosdbInstanceGrp

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbClusterInstanceGroupsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceGrpInfoList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceGrpInfoList))
	tmpList := make([]map[string]interface{}, 0, len(instanceGrpInfoList))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}

/*
Use this data source to query detailed information of tsf group_instances

Example Usage

```hcl
data "tencentcloud_tsf_group_instances" "group_instances" {
  group_id = "group-yrjkln9v"
  search_word = "testing"
  order_by = "ASC"
  order_type = 0
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfGroupInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfGroupInstancesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "group id.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search word.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "order term.",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "order type.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Machine information of the deployment group.Note: This field may return null, which means no valid value was found.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of machine instances.Note: This field may return null, which means no valid value was found.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of machine instances.Note: This field may return null, which means no valid value was found.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Machine instance ID.Note: This field may return null, which means no valid value was found.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Machine name.Note: This field may return null, which means no valid value was found.",
									},
									"lan_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Private IP address.Note: This field may return null, which means no valid value was found.",
									},
									"wan_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Public IP address.Note: This field may return null, which means no valid value was found.",
									},
									"instance_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description.Note: This field may return null, which means no valid value was found.",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster id.Note: This field may return null, which means no valid value was found.",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster name. Note: This field may return null, which means no valid value was found.",
									},
									"instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VM status. For virtual machines, it indicates the status of the virtual machine. For containers, it indicates the status of the virtual machine where the pod is located.Note: This field may return null, which means no valid value was found.",
									},
									"instance_available_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VM availability status. For virtual machines, it indicates whether the virtual machine can be used as a resource. For containers, it indicates whether the virtual machine can be used to deploy pods.Note: This field may return null, which means no valid value was found.",
									},
									"service_instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of service instances under the service. For virtual machines, it indicates whether the application is available and the agent status. For containers, it indicates the status of the pod.Note: This field may return null, which means no valid value was found.",
									},
									"count_in_tsf": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Indicates whether this instance has been added to the TSF.Note: This field may return null, which means no valid value was found.",
									},
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group id.Note: This field may return null, which means no valid value was found.",
									},
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application id.Note: This field may return null, which means no valid value was found.",
									},
									"application_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application name. Note: This field may return null, which means no valid value was found.",
									},
									"instance_created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creation time of the machine instance in CVM.Note: This field may return null, which means no valid value was found.",
									},
									"instance_expired_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire time of the machine instance in CVM.Note: This field may return null, which means no valid value was found.",
									},
									"instance_charge_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "machine instance charge type.Note: This field may return null, which means no valid value was found.",
									},
									"instance_total_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Total CPU information of the machine instance.Note: This field may return null, which means no valid value was found.",
									},
									"instance_total_mem": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Total memory information of the machine instance.Note: This field may return null, which means no valid value was found.",
									},
									"instance_used_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "CPU information used by the machine instance.Note: This field may return null, which means no valid value was found.",
									},
									"instance_used_mem": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Memory information used by the machine instance.Note: This field may return null, which means no valid value was found.",
									},
									"instance_limit_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Limit CPU information of the machine instance.Note: This field may return null, which means no valid value was found.",
									},
									"instance_limit_mem": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Limit memory information of the machine instance.Note: This field may return null, which means no valid value was found.",
									},
									"instance_pkg_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance pkg version.Note: This field may return null, which means no valid value was found.",
									},
									"cluster_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster type.Note: This field may return null, which means no valid value was found.",
									},
									"restrict_state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Business status of the machine instance.Note: This field may return null, which means no valid value was found.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Update time.Note: This field may return null, which means no valid value was found.",
									},
									"operation_state": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Execution status of the instance.Note: This field may return null, which means no valid value was found.",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace id.Note: This field may return null, which means no valid value was found.",
									},
									"instance_zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance zone id.Note: This field may return null, which means no valid value was found.",
									},
									"instance_import_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "InstanceImportMode import mode.Note: This field may return null, which means no valid value was found.",
									},
									"application_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application id.Note: This field may return null, which means no valid value was found.",
									},
									"application_resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "application resource id.Note: This field may return null, which means no valid value was found.",
									},
									"service_sidecar_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Sidecar status.Note: This field may return null, which means no valid value was found.",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group name.Note: This field may return null, which means no valid value was found.",
									},
									"namespace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace name.Note: This field may return null, which means no valid value was found.",
									},
									"reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Health checking reason.Note: This field may return null, which means no valid value was found.",
									},
									"agent_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Agent version.Note: This field may return null, which means no valid value was found.",
									},
									"node_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Container host instance ID.Note: This field may return null, which means no valid value was found.",
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

func dataSourceTencentCloudTsfGroupInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_group_instances.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("order_type"); v != nil {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *tsf.TsfPageInstance
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTsfGroupInstancesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result.Content))
	tsfPageInstanceMap := map[string]interface{}{}
	if result != nil {
		if result.TotalCount != nil {
			tsfPageInstanceMap["total_count"] = result.TotalCount
		}

		if result.Content != nil {
			contentList := []interface{}{}
			for _, content := range result.Content {
				contentMap := map[string]interface{}{}

				if content.InstanceId != nil {
					contentMap["instance_id"] = content.InstanceId
				}

				if content.InstanceName != nil {
					contentMap["instance_name"] = content.InstanceName
				}

				if content.LanIp != nil {
					contentMap["lan_ip"] = content.LanIp
				}

				if content.WanIp != nil {
					contentMap["wan_ip"] = content.WanIp
				}

				if content.InstanceDesc != nil {
					contentMap["instance_desc"] = content.InstanceDesc
				}

				if content.ClusterId != nil {
					contentMap["cluster_id"] = content.ClusterId
				}

				if content.ClusterName != nil {
					contentMap["cluster_name"] = content.ClusterName
				}

				if content.InstanceStatus != nil {
					contentMap["instance_status"] = content.InstanceStatus
				}

				if content.InstanceAvailableStatus != nil {
					contentMap["instance_available_status"] = content.InstanceAvailableStatus
				}

				if content.ServiceInstanceStatus != nil {
					contentMap["service_instance_status"] = content.ServiceInstanceStatus
				}

				if content.CountInTsf != nil {
					contentMap["count_in_tsf"] = content.CountInTsf
				}

				if content.GroupId != nil {
					contentMap["group_id"] = content.GroupId
				}

				if content.ApplicationId != nil {
					contentMap["application_id"] = content.ApplicationId
				}

				if content.ApplicationName != nil {
					contentMap["application_name"] = content.ApplicationName
				}

				if content.InstanceCreatedTime != nil {
					contentMap["instance_created_time"] = content.InstanceCreatedTime
				}

				if content.InstanceExpiredTime != nil {
					contentMap["instance_expired_time"] = content.InstanceExpiredTime
				}

				if content.InstanceChargeType != nil {
					contentMap["instance_charge_type"] = content.InstanceChargeType
				}

				if content.InstanceTotalCpu != nil {
					contentMap["instance_total_cpu"] = content.InstanceTotalCpu
				}

				if content.InstanceTotalMem != nil {
					contentMap["instance_total_mem"] = content.InstanceTotalMem
				}

				if content.InstanceUsedCpu != nil {
					contentMap["instance_used_cpu"] = content.InstanceUsedCpu
				}

				if content.InstanceUsedMem != nil {
					contentMap["instance_used_mem"] = content.InstanceUsedMem
				}

				if content.InstanceLimitCpu != nil {
					contentMap["instance_limit_cpu"] = content.InstanceLimitCpu
				}

				if content.InstanceLimitMem != nil {
					contentMap["instance_limit_mem"] = content.InstanceLimitMem
				}

				if content.InstancePkgVersion != nil {
					contentMap["instance_pkg_version"] = content.InstancePkgVersion
				}

				if content.ClusterType != nil {
					contentMap["cluster_type"] = content.ClusterType
				}

				if content.RestrictState != nil {
					contentMap["restrict_state"] = content.RestrictState
				}

				if content.UpdateTime != nil {
					contentMap["update_time"] = content.UpdateTime
				}

				if content.OperationState != nil {
					contentMap["operation_state"] = content.OperationState
				}

				if content.NamespaceId != nil {
					contentMap["namespace_id"] = content.NamespaceId
				}

				if content.InstanceZoneId != nil {
					contentMap["instance_zone_id"] = content.InstanceZoneId
				}

				if content.InstanceImportMode != nil {
					contentMap["instance_import_mode"] = content.InstanceImportMode
				}

				if content.ApplicationType != nil {
					contentMap["application_type"] = content.ApplicationType
				}

				if content.ApplicationResourceType != nil {
					contentMap["application_resource_type"] = content.ApplicationResourceType
				}

				if content.ServiceSidecarStatus != nil {
					contentMap["service_sidecar_status"] = content.ServiceSidecarStatus
				}

				if content.GroupName != nil {
					contentMap["group_name"] = content.GroupName
				}

				if content.NamespaceName != nil {
					contentMap["namespace_name"] = content.NamespaceName
				}

				if content.Reason != nil {
					contentMap["reason"] = content.Reason
				}

				if content.AgentVersion != nil {
					contentMap["agent_version"] = content.AgentVersion
				}

				if content.NodeInstanceId != nil {
					contentMap["node_instance_id"] = content.NodeInstanceId
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.InstanceId)
			}

			tsfPageInstanceMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageInstanceMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tsfPageInstanceMap); e != nil {
			return e
		}
	}
	return nil
}

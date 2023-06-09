/*
Use this data source to query detailed information of tsf pod_instances

Example Usage

```hcl
data "tencentcloud_tsf_pod_instances" "pod_instances" {
  group_id = "group-ynd95rea"
  pod_name_list = ["keep-terraform-6f8f977688-zvphm"]
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

func dataSourceTencentCloudTsfPodInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfPodInstancesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance&amp;#39;s group ID.",
			},

			"pod_name_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter, pod name list.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "pod instance list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of records.Note: This field may return null, which means no valid value was found.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Content list.Note: This field may return null, which means no valid value was found.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pod_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance name (corresponding to the pod name in Kubernetes).Note: This field may return null, which means no valid value was found.",
									},
									"pod_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance id (corresponding to the pod instance id in Kubernetes).Note: This field may return null, which means no valid value was found.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance status. Please refer to the definition of instance and container status below. Starting (pod not ready): Starting; Running: Running; Abnormal: Abnormal; Stopped: Stopped;Note: This field may return null, which means no valid value was found.",
									},
									"reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance reason for current status.Note: This field may return null, which means no valid value was found.",
									},
									"node_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance node ip.Note: This field may return null, which means no valid value was found.",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ip.Note: This field may return null, which means no valid value was found.",
									},
									"restart_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance restart count.Note: This field may return null, which means no valid value was found.",
									},
									"ready_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance ready count.Note: This field may return null, which means no valid value was found.",
									},
									"runtime": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance run time.Note: This field may return null, which means no valid value was found.",
									},
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance start time.Note: This field may return null, which means no valid value was found.",
									},
									"service_instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance serve status.Note: This field may return null, which means no valid value was found.",
									},
									"instance_available_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance available status.Note: This field may return null, which means no valid value was found.",
									},
									"instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance status.Note: This field may return null, which means no valid value was found.",
									},
									"node_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance node id.Note: This field may return null, which means no valid value was found.",
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

func dataSourceTencentCloudTsfPodInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_pod_instances.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pod_name_list"); ok {
		podNameListSet := v.(*schema.Set).List()
		paramMap["PodNameList"] = helper.InterfacesStringsPoint(podNameListSet)
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	var groupPodResult *tsf.GroupPodResult
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfDescribePodInstancesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		groupPodResult = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(groupPodResult.Content))
	groupPodResultMap := map[string]interface{}{}
	if groupPodResult != nil {
		if groupPodResult.TotalCount != nil {
			groupPodResultMap["total_count"] = groupPodResult.TotalCount
		}

		if groupPodResult.Content != nil {
			contentList := []interface{}{}
			for _, content := range groupPodResult.Content {
				contentMap := map[string]interface{}{}

				if content.PodName != nil {
					contentMap["pod_name"] = content.PodName
				}

				if content.PodId != nil {
					contentMap["pod_id"] = content.PodId
				}

				if content.Status != nil {
					contentMap["status"] = content.Status
				}

				if content.Reason != nil {
					contentMap["reason"] = content.Reason
				}

				if content.NodeIp != nil {
					contentMap["node_ip"] = content.NodeIp
				}

				if content.Ip != nil {
					contentMap["ip"] = content.Ip
				}

				if content.RestartCount != nil {
					contentMap["restart_count"] = content.RestartCount
				}

				if content.ReadyCount != nil {
					contentMap["ready_count"] = content.ReadyCount
				}

				if content.Runtime != nil {
					contentMap["runtime"] = content.Runtime
				}

				if content.CreatedAt != nil {
					contentMap["created_at"] = content.CreatedAt
				}

				if content.ServiceInstanceStatus != nil {
					contentMap["service_instance_status"] = content.ServiceInstanceStatus
				}

				if content.InstanceAvailableStatus != nil {
					contentMap["instance_available_status"] = content.InstanceAvailableStatus
				}

				if content.InstanceStatus != nil {
					contentMap["instance_status"] = content.InstanceStatus
				}

				if content.NodeInstanceId != nil {
					contentMap["node_instance_id"] = content.NodeInstanceId
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.PodName)
			}

			groupPodResultMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{groupPodResultMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), groupPodResultMap); e != nil {
			return e
		}
	}
	return nil
}

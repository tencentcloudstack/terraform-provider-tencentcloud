/*
Use this data source to query detailed information of tsf container_group

Example Usage

```hcl
data "tencentcloud_tsf_container_group" "container_group" {
  application_id = "application-a24x29xv"
  search_word = "keep"
  order_by = "createTime"
  order_type = 0
  cluster_id = "cluster-vwgj5e6y"
  namespace_id = "namespace-aemrg36v"
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

func dataSourceTencentCloudTsfContainerGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfContainerGroupRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search word, support group name.",
			},

			"application_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ApplicationId, required.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The sorting field. By default, it is the createTime field. Supports id, name, createTime.",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The sorting order. By default, it is 1, indicating descending order. 0 indicates ascending order, and 1 indicates descending order.",
			},

			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster Id.",
			},

			"namespace_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Namespace Id.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "result list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of deployment groups.Note: This field may return null, indicating that no valid value was found.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group Id.Note: This field may return null, indicating that no valid value was found.",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group name.Note: This field may return null, indicating that no valid value was found.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Create time.Note: This field may return null, indicating that no valid value was found.",
									},
									"server": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image server.Note: This field may return null, indicating that no valid value was found.",
									},
									"repo_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image name.Note: This field may return null, indicating that no valid value was found.",
									},
									"tag_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image version Name.Note: This field may return null, indicating that no valid value was found.",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster Id.Note: This field may return null, indicating that no valid value was found.",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster name.Note: This field may return null, indicating that no valid value was found.",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace Id.Note: This field may return null, indicating that no valid value was found.",
									},
									"namespace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace name.Note: This field may return null, indicating that no valid value was found.",
									},
									"cpu_request": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The initial amount of CPU, corresponding to K8S request.Note: This field may return null, indicating that no valid value was found.",
									},
									"cpu_limit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The maximum amount of CPU, corresponding to K8S limit.Note: This field may return null, indicating that no valid value was found.",
									},
									"mem_request": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The initial amount of memory allocated in MiB, corresponding to K8S request.Note: This field may return null, indicating that no valid value was found.",
									},
									"mem_limit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The maximum amount of memory allocated in MiB, corresponding to K8S limit.Note: This field may return null, indicating that no valid value was found.",
									},
									"alias": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Group description.Note: This field may return null, indicating that no valid value was found.",
									},
									"kube_inject_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The value of KubeInjectEnable.Note: This field may return null, indicating that no valid value was found.",
									},
									"updated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Update type.Note: This field may return null, indicating that no valid value was found.",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total count.",
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

func dataSourceTencentCloudTsfContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf__container_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_id"); ok {
		paramMap["ApplicationId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("order_type"); v != nil {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		paramMap["NamespaceId"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *tsf.ContainGroupResult
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTsfDescriptionContainerGroupByFilter(ctx, paramMap)
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
	containGroupResultMap := map[string]interface{}{}
	if result != nil {

		if result.Content != nil {
			contentList := []interface{}{}
			for _, content := range result.Content {
				contentMap := map[string]interface{}{}

				if content.GroupId != nil {
					contentMap["group_id"] = content.GroupId
				}

				if content.GroupName != nil {
					contentMap["group_name"] = content.GroupName
				}

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.Server != nil {
					contentMap["server"] = content.Server
				}

				if content.RepoName != nil {
					contentMap["repo_name"] = content.RepoName
				}

				if content.TagName != nil {
					contentMap["tag_name"] = content.TagName
				}

				if content.ClusterId != nil {
					contentMap["cluster_id"] = content.ClusterId
				}

				if content.ClusterName != nil {
					contentMap["cluster_name"] = content.ClusterName
				}

				if content.NamespaceId != nil {
					contentMap["namespace_id"] = content.NamespaceId
				}

				if content.NamespaceName != nil {
					contentMap["namespace_name"] = content.NamespaceName
				}

				if content.CpuRequest != nil {
					contentMap["cpu_request"] = content.CpuRequest
				}

				if content.CpuLimit != nil {
					contentMap["cpu_limit"] = content.CpuLimit
				}

				if content.MemRequest != nil {
					contentMap["mem_request"] = content.MemRequest
				}

				if content.MemLimit != nil {
					contentMap["mem_limit"] = content.MemLimit
				}

				if content.Alias != nil {
					contentMap["alias"] = content.Alias
				}

				if content.KubeInjectEnable != nil {
					contentMap["kube_inject_enable"] = content.KubeInjectEnable
				}

				if content.UpdatedTime != nil {
					contentMap["updated_time"] = content.UpdatedTime
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.GroupId)
			}

			containGroupResultMap["content"] = contentList
		}

		if result.TotalCount != nil {
			containGroupResultMap["total_count"] = result.TotalCount
		}

		_ = d.Set("result", []interface{}{containGroupResultMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), containGroupResultMap); e != nil {
			return e
		}
	}
	return nil
}

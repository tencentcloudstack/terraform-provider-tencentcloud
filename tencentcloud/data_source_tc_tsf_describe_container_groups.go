/*
Use this data source to query detailed information of tsf describe_container_groups

Example Usage

```hcl
data "tencentcloud_tsf_describe_container_groups" "describe_container_groups" {
  search_word = ""
  application_id = ""
  order_by = ""
  order_type =
  cluster_id = ""
  namespace_id = ""
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

func dataSourceTencentCloudTsfDescribeContainerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfDescribeContainerGroupsRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search field, fuzzy search the groupName field.",
			},

			"application_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The ID of the application to which the group belongs. .",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field, default is createTime field. Supports id, name, createTime.",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Order type. Pass 0 for ascending order and 1 for descending order.",
			},

			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The cluster Id of  group.",
			},

			"namespace_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Namespace ID .",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The result.",
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
										Description: " group ID.Note: This field may return null, indicating that no valid value was found.",
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
										Description: "Image name, such as /tsf/nginx.Note: This field may return null, indicating that no valid value was found.",
									},
									"tag_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image tag.Note: This field may return null, indicating that no valid value was found.",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ClusterId.Note: This field may return null, indicating that no valid value was found.",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ClusterName.Note: This field may return null, indicating that no valid value was found.",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace Id.Note: This field may return null, indicating that no valid value was found.",
									},
									"namespace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace Name.Note: This field may return null, indicating that no valid value was found.",
									},
									"cpu_request": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cpu Request.Note: This field may return null, indicating that no valid value was found.",
									},
									"cpu_limit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cpu Limit.Note: This field may return null, indicating that no valid value was found.",
									},
									"mem_request": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Memory request.Note: This field may return null, indicating that no valid value was found.",
									},
									"mem_limit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Memory Limit.Note: This field may return null, indicating that no valid value was found.",
									},
									"alias": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy Group alias.Note: This field may return null, indicating that no valid value was found.",
									},
									"kube_inject_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether import form TKE or not.Note: This field may return null, indicating that no valid value was found.",
									},
									"updated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Update time.Note: This field may return null, indicating that no valid value was found.",
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

func dataSourceTencentCloudTsfDescribeContainerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_describe_container_groups.read")()
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

	var result []*tsf.ContainGroupResult

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfDescribeContainerGroupsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		containGroupResultMap := map[string]interface{}{}

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
			}

			containGroupResultMap["content"] = []interface{}{contentList}
		}

		if result.TotalCount != nil {
			containGroupResultMap["total_count"] = result.TotalCount
		}

		ids = append(ids, *result.ClusterId)
		_ = d.Set("result", containGroupResultMap)
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

package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfGroupsRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "searchWord, support groupName.",
			},

			"application_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "applicationId.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "sort term.",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "order type, 0 desc, 1 asc.",
			},

			"namespace_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "namespace Id.",
			},

			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "clusterId.",
			},

			"group_resource_type_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Group resourceType list.",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "group status filter, `Running`: running, `Unknown`: unknown.",
			},

			"group_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "group Id list.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Pagination information of the virtual machine deployment group.Note: This field may return null, indicating that no valid value was found.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total count virtual machine deployment group. Note: This field may return null, indicating that no valid value was found.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Virtual machine deployment group list. Note: This field may return null, indicating that no valid value was found.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group ID. Note: This field may return null, indicating that no valid value was found.",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group ID. Note: This field may return null, indicating that no valid value was found.",
									},
									"application_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application type. Note: This field may return null, indicating that no valid value was found.",
									},
									"group_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group description. Note: This field may return null, indicating that no valid value was found.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group update time. Note: This field may return null, indicating that no valid value was found.",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster ID. Note: This field may return null, indicating that no valid value was found.",
									},
									"startup_parameters": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group start up Parameters. Note: This field may return null, indicating that no valid value was found.",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace ID. Note: This field may return null, indicating that no valid value was found.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Create Time. Note: This field may return null, indicating that no valid value was found.",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster name. Note: This field may return null, indicating that no valid value was found.",
									},
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application ID. Note: This field may return null, indicating that no valid value was found.",
									},
									"application_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application name. Note: This field may return null, indicating that no valid value was found.",
									},
									"namespace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace name. Note: This field may return null, indicating that no valid value was found.",
									},
									"microservice_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice type. Note: This field may return null, indicating that no valid value was found.",
									},
									"group_resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group resource type. Note: This field may return null, indicating that no valid value was found.",
									},
									"updated_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Update time. Note: This field may return null, indicating that no valid value was found.",
									},
									"deploy_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group description. Note: This field may return null, indicating that no valid value was found.",
									},
									"alias": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group alias. Note: This field may return null, indicating that no valid value was found.",
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

func dataSourceTencentCloudTsfGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_groups.read")()
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

	if v, ok := d.GetOk("namespace_id"); ok {
		paramMap["NamespaceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_resource_type_list"); ok {
		groupResourceTypeListSet := v.(*schema.Set).List()
		paramMap["GroupResourceTypeList"] = helper.InterfacesStringsPoint(groupResourceTypeListSet)
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id_list"); ok {
		groupIdListSet := v.(*schema.Set).List()
		paramMap["GroupIdList"] = helper.InterfacesStringsPoint(groupIdListSet)
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *tsf.TsfPageVmGroup
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTsfGroupsByFilter(ctx, paramMap)
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
	tsfPageVmGroupMap := map[string]interface{}{}
	if result != nil {
		if result.TotalCount != nil {
			tsfPageVmGroupMap["total_count"] = result.TotalCount
		}

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

				if content.ApplicationType != nil {
					contentMap["application_type"] = content.ApplicationType
				}

				if content.GroupDesc != nil {
					contentMap["group_desc"] = content.GroupDesc
				}

				if content.UpdateTime != nil {
					contentMap["update_time"] = content.UpdateTime
				}

				if content.ClusterId != nil {
					contentMap["cluster_id"] = content.ClusterId
				}

				if content.StartupParameters != nil {
					contentMap["startup_parameters"] = content.StartupParameters
				}

				if content.NamespaceId != nil {
					contentMap["namespace_id"] = content.NamespaceId
				}

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.ClusterName != nil {
					contentMap["cluster_name"] = content.ClusterName
				}

				if content.ApplicationId != nil {
					contentMap["application_id"] = content.ApplicationId
				}

				if content.ApplicationName != nil {
					contentMap["application_name"] = content.ApplicationName
				}

				if content.NamespaceName != nil {
					contentMap["namespace_name"] = content.NamespaceName
				}

				if content.MicroserviceType != nil {
					contentMap["microservice_type"] = content.MicroserviceType
				}

				if content.GroupResourceType != nil {
					contentMap["group_resource_type"] = content.GroupResourceType
				}

				if content.UpdatedTime != nil {
					contentMap["updated_time"] = content.UpdatedTime
				}

				if content.DeployDesc != nil {
					contentMap["deploy_desc"] = content.DeployDesc
				}

				if content.Alias != nil {
					contentMap["alias"] = content.Alias
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.GroupId)
			}

			tsfPageVmGroupMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageVmGroupMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tsfPageVmGroupMap); e != nil {
			return e
		}
	}
	return nil
}

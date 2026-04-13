package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudConfigDiscoveredResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudConfigDiscoveredResourcesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. Supported filter names: resourceName (resource name), resourceId (resource ID).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name. Valid values: resourceName, resourceId.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Filter field values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tag filter conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"order_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort type. Valid values: asc, desc.",
			},

			"resource_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Discovered resource list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type.",
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource name.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID.",
						},
						"resource_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource region.",
						},
						"resource_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource status.",
						},
						"resource_delete": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource deletion mark. Valid values: 1 (deleted), 2 (not deleted).",
						},
						"resource_create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource creation time.",
						},
						"resource_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource availability zone.",
						},
						"compliance_result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance result. Valid values: COMPLIANT, NON_COMPLIANT.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Resource tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
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

func dataSourceTencentCloudConfigDiscoveredResourcesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_config_discovered_resources.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := buildDiscoveredResourcesParamMap(d)

	respData, reqErr := service.DescribeConfigDiscoveredResourcesByFilter(ctx, paramMap)
	if reqErr != nil {
		return reqErr
	}

	resourceList := flattenResourceListInfo(respData)
	_ = d.Set("resource_list", resourceList)

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

func buildDiscoveredResourcesParamMap(d *schema.ResourceData) map[string]interface{} {
	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("filters"); ok {
		rawList := v.([]interface{})
		filters := make([]*configv20220802.Filter, 0, len(rawList))
		for _, item := range rawList {
			filterMap := item.(map[string]interface{})
			filter := &configv20220802.Filter{}
			if name, ok := filterMap["name"].(string); ok && name != "" {
				filter.Name = helper.String(name)
			}

			if vals, ok := filterMap["values"]; ok {
				valList := vals.([]interface{})
				for _, v := range valList {
					val := v.(string)
					filter.Values = append(filter.Values, &val)
				}
			}

			filters = append(filters, filter)
		}

		paramMap["Filters"] = filters
	}

	if v, ok := d.GetOk("tags"); ok {
		rawList := v.([]interface{})
		tags := make([]*configv20220802.Tag, 0, len(rawList))
		for _, item := range rawList {
			tagMap := item.(map[string]interface{})
			tag := &configv20220802.Tag{}
			if k, ok := tagMap["tag_key"].(string); ok && k != "" {
				tag.TagKey = helper.String(k)
			}

			if val, ok := tagMap["tag_value"].(string); ok && val != "" {
				tag.TagValue = helper.String(val)
			}

			tags = append(tags, tag)
		}

		paramMap["Tags"] = tags
	}

	if v, ok := d.GetOk("order_type"); ok {
		paramMap["OrderType"] = v.(string)
	}

	return paramMap
}

func flattenResourceListInfo(items []*configv20220802.ResourceListInfo) []map[string]interface{} {
	resourceList := make([]map[string]interface{}, 0, len(items))
	for _, res := range items {
		resMap := map[string]interface{}{}

		if res.ResourceType != nil {
			resMap["resource_type"] = res.ResourceType
		}

		if res.ResourceName != nil {
			resMap["resource_name"] = res.ResourceName
		}

		if res.ResourceId != nil {
			resMap["resource_id"] = res.ResourceId
		}

		if res.ResourceRegion != nil {
			resMap["resource_region"] = res.ResourceRegion
		}

		if res.ResourceStatus != nil {
			resMap["resource_status"] = res.ResourceStatus
		}

		if res.ResourceDelete != nil {
			resMap["resource_delete"] = int(*res.ResourceDelete)
		}

		if res.ResourceCreateTime != nil {
			resMap["resource_create_time"] = res.ResourceCreateTime
		}

		if res.ResourceZone != nil {
			resMap["resource_zone"] = res.ResourceZone
		}

		if res.ComplianceResult != nil {
			resMap["compliance_result"] = res.ComplianceResult
		}

		if res.Tags != nil {
			tagList := make([]map[string]interface{}, 0, len(res.Tags))
			for _, tag := range res.Tags {
				tagMap := map[string]interface{}{}
				if tag.TagKey != nil {
					tagMap["tag_key"] = tag.TagKey
				}

				if tag.TagValue != nil {
					tagMap["tag_value"] = tag.TagValue
				}

				tagList = append(tagList, tagMap)
			}

			resMap["tags"] = tagList
		}

		resourceList = append(resourceList, resMap)
	}

	return resourceList
}

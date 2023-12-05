package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudGaapProxyGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapProxyGroupsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Project ID. Value range:-1, All projects under this user0, default projectOther values, specified items.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions,The upper limit of Filter.Values per request is 5.RealServerRegion - String - Required: No - (filtering criteria) Filter by real server region, refer to the RegionId in the returned results of the DescribeDestRegions interface.PackageType - String - Required: No - (Filter condition) proxy group type, where &amp;#39;Thunder&amp;#39; represents the standard proxy group and &amp;#39;Accelerator&amp;#39; represents the silver acceleration proxy group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter conditions.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "filtering value.",
						},
					},
				},
			},

			"tag_set": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Tag list, when this field exists, pulls the resource list under the corresponding tag.Supports a maximum of 5 labels. When there are two or more labels and any one of them is met, the proxy group will be pulled out.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag Key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag Value.",
						},
					},
				},
			},

			"proxy_group_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of proxy groups.Note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy group Id.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy group domain nameNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy Group NameNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project Id.",
						},
						"real_server_region_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Real Server Region Info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region Id.",
									},
									"region_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region Name.",
									},
									"region_area": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region of the computer room.",
									},
									"region_area_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region name of the computer room.",
									},
									"idc_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of computer room, where &#39;dc&#39; represents the DataCenter data center and &#39;ec&#39; represents the EdgeComputing edge node.",
									},
									"feature_bitmap": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Property bitmap, where each bit represents a property, where:0, indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"support_feature": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Ability to access regional supportNote: This field may return null, indicating that a valid value cannot be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"network_type": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "A list of network types supported by the access area, with &#39;normal&#39; indicating support for regular BGP, &#39;cn2&#39; indicating premium BGP, &#39;triple&#39; indicating three networks, and &#39;secure_EIP&#39; represents a custom secure EIP.",
												},
											},
										},
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy group status.Among them,&#39;RUNNING&#39; indicates running;&#39;CREATING&#39; indicates being created;&#39;DESTROYING&#39; indicates being destroyed;&#39;MOVING&#39; indicates that the proxy is being migrated;&#39;CHANGING&#39; indicates partial deployment.",
						},
						"tag_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag Set.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag Key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag Value.",
									},
								},
							},
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy Group VersionNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Create TimeNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"proxy_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Does the proxy group include Microsoft proxysNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"http3_supported": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Supports the identification of Http3 features, where:0 indicates shutdown;1 indicates enabled.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"feature_bitmap": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Property bitmap, where each bit represents a property, where:0, indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudGaapProxyGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_proxy_groups.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*gaap.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := gaap.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	if v, ok := d.GetOk("tag_set"); ok {
		tagSetSet := v.([]interface{})
		tmpSet := make([]*gaap.TagPair, 0, len(tagSetSet))

		for _, item := range tagSetSet {
			tagPair := gaap.TagPair{}
			tagPairMap := item.(map[string]interface{})

			if v, ok := tagPairMap["tag_key"]; ok {
				tagPair.TagKey = helper.String(v.(string))
			}
			if v, ok := tagPairMap["tag_value"]; ok {
				tagPair.TagValue = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &tagPair)
		}
		paramMap["tag_set"] = tmpSet
	}

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

	var proxyGroupList []*gaap.ProxyGroupInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapProxyGroupsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		proxyGroupList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(proxyGroupList))
	tmpList := make([]map[string]interface{}, 0, len(proxyGroupList))

	if proxyGroupList != nil {
		for _, proxyGroupInfo := range proxyGroupList {
			proxyGroupInfoMap := map[string]interface{}{}

			if proxyGroupInfo.GroupId != nil {
				proxyGroupInfoMap["group_id"] = proxyGroupInfo.GroupId
			}

			if proxyGroupInfo.Domain != nil {
				proxyGroupInfoMap["domain"] = proxyGroupInfo.Domain
			}

			if proxyGroupInfo.GroupName != nil {
				proxyGroupInfoMap["group_name"] = proxyGroupInfo.GroupName
			}

			if proxyGroupInfo.ProjectId != nil {
				proxyGroupInfoMap["project_id"] = proxyGroupInfo.ProjectId
			}

			if proxyGroupInfo.RealServerRegionInfo != nil {
				realServerRegionInfoMap := map[string]interface{}{}

				if proxyGroupInfo.RealServerRegionInfo.RegionId != nil {
					realServerRegionInfoMap["region_id"] = proxyGroupInfo.RealServerRegionInfo.RegionId
				}

				if proxyGroupInfo.RealServerRegionInfo.RegionName != nil {
					realServerRegionInfoMap["region_name"] = proxyGroupInfo.RealServerRegionInfo.RegionName
				}

				if proxyGroupInfo.RealServerRegionInfo.RegionArea != nil {
					realServerRegionInfoMap["region_area"] = proxyGroupInfo.RealServerRegionInfo.RegionArea
				}

				if proxyGroupInfo.RealServerRegionInfo.RegionAreaName != nil {
					realServerRegionInfoMap["region_area_name"] = proxyGroupInfo.RealServerRegionInfo.RegionAreaName
				}

				if proxyGroupInfo.RealServerRegionInfo.IDCType != nil {
					realServerRegionInfoMap["idc_type"] = proxyGroupInfo.RealServerRegionInfo.IDCType
				}

				if proxyGroupInfo.RealServerRegionInfo.FeatureBitmap != nil {
					realServerRegionInfoMap["feature_bitmap"] = proxyGroupInfo.RealServerRegionInfo.FeatureBitmap
				}

				if proxyGroupInfo.RealServerRegionInfo.SupportFeature != nil {
					supportFeatureMap := map[string]interface{}{}

					if proxyGroupInfo.RealServerRegionInfo.SupportFeature.NetworkType != nil {
						supportFeatureMap["network_type"] = proxyGroupInfo.RealServerRegionInfo.SupportFeature.NetworkType
					}

					realServerRegionInfoMap["support_feature"] = []interface{}{supportFeatureMap}
				}

				proxyGroupInfoMap["real_server_region_info"] = []interface{}{realServerRegionInfoMap}
			}

			if proxyGroupInfo.Status != nil {
				proxyGroupInfoMap["status"] = proxyGroupInfo.Status
			}

			if proxyGroupInfo.TagSet != nil {
				tagSetList := []interface{}{}
				for _, tagSet := range proxyGroupInfo.TagSet {
					tagSetMap := map[string]interface{}{}

					if tagSet.TagKey != nil {
						tagSetMap["tag_key"] = tagSet.TagKey
					}

					if tagSet.TagValue != nil {
						tagSetMap["tag_value"] = tagSet.TagValue
					}

					tagSetList = append(tagSetList, tagSetMap)
				}

				proxyGroupInfoMap["tag_set"] = tagSetList
			}

			if proxyGroupInfo.Version != nil {
				proxyGroupInfoMap["version"] = proxyGroupInfo.Version
			}

			if proxyGroupInfo.CreateTime != nil {
				proxyGroupInfoMap["create_time"] = proxyGroupInfo.CreateTime
			}

			if proxyGroupInfo.ProxyType != nil {
				proxyGroupInfoMap["proxy_type"] = proxyGroupInfo.ProxyType
			}

			if proxyGroupInfo.Http3Supported != nil {
				proxyGroupInfoMap["http3_supported"] = proxyGroupInfo.Http3Supported
			}

			if proxyGroupInfo.FeatureBitmap != nil {
				proxyGroupInfoMap["feature_bitmap"] = proxyGroupInfo.FeatureBitmap
			}

			ids = append(ids, *proxyGroupInfo.GroupId)
			tmpList = append(tmpList, proxyGroupInfoMap)
		}

		_ = d.Set("proxy_group_list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}

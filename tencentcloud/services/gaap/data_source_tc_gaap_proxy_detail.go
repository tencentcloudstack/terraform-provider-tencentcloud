package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
)

func DataSourceTencentCloudGaapProxyDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapProxyDetailRead,
		Schema: map[string]*schema.Schema{
			"proxy_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Proxy Id.",
			},

			"proxy_detail": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Proxy Detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "(Old parameter, please use ProxyId) Proxy instance ID.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The creation time, using a Unix timestamp, represents the number of seconds that have passed since January 1, 1970 (midnight UTC/GMT).",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project Id.",
						},
						"proxy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Proxy Name.",
						},
						"access_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access Region.",
						},
						"real_server_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real Server Region.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Band width, in Mbps.",
						},
						"concurrent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Concurrent, in 10000 pieces/second.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy status. Among them:RUNNING indicates running;CREATING indicates being created;DESTROYING indicates being destroyed;OPENING indicates being opened;CLOSING indicates being closed;Closed indicates that it has been closed;ADJUSTING represents a configuration change in progress;ISOLATING indicates being isolated;ISOLATED indicates that it has been isolated;CLONING indicates copying;RECOVERING indicates that the proxy is being maintained;MOVING indicates that migration is in progress.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version 1.0, 2.0, 3.0.",
						},
						"proxy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "(New parameter) proxy instance ID.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"scalarable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "1. This proxy can be scaled and expanded; 0, this proxy cannot be scaled or expanded.",
						},
						"support_protocols": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Supported protocol types.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy group ID, which exists when a proxy belongs to a certain proxy group.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security policy ID, which exists when a security policy is set.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"access_region_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detailed information about the access region, including the region ID and domain name.Note: This field may return null, indicating that a valid value cannot be obtained.",
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
										Description: "The type of computer room, where dc represents the DataCenter data center and ec represents the EdgeComputing edge node.",
									},
									"feature_bitmap": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Property bitmap, where each bit represents a property, where:0 indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.",
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
													Description: "A list of network types supported by the access area, with normal indicating support for regular BGP, cn2 indicating premium BGP, triple indicating three networks, and secure_EIP represents a custom secure EIP.",
												},
											},
										},
									},
								},
							},
						},
						"real_server_region_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detailed information of the real server region, including the region ID and domain name.Note: This field may return null, indicating that a valid value cannot be obtained.",
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
										Description: "The type of computer room, where dc represents the DataCenter data center and ec represents the EdgeComputing edge node.",
									},
									"feature_bitmap": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Property bitmap, where each bit represents a property, where:0 indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.",
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
													Description: "A list of network types supported by the access area, with normal indicating support for regular BGP, cn2 indicating premium BGP, triple indicating three networks, and secure_EIP represents a custom secure EIP.",
												},
											},
										},
									},
								},
							},
						},
						"forward_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy forwarding IP.",
						},
						"tag_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "tag list, when there are no labels, this field is an empty list.Note: This field may return null, indicating that a valid value cannot be obtained.",
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
						"support_security": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Does it support security group configurationNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"billing_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing type: 0 represents bandwidth based billing, and 1 represents traffic based billing.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"related_global_domains": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of domain names associated with resolutionNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"modify_config_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Configuration change timeNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"proxy_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "proxy type, 100 represents THUNDER proxy, 103 represents Microsoft cooperation proxyNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"client_ip_method": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "The method of obtaining client IP through proxys, where 0 represents TOA and 1 represents Proxy ProtocolNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"ip_address_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version: IPv4, IPv6Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network type: normal represents regular BGP, cn2 represents premium BGP, triple represents triple network, secure_EIP represents customized security EIPNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"package_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy package type: Thunder represents standard proxy, Accelerator represents silver acceleration proxy,CrossBorder represents a cross-border proxy.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"ban_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Blocking and Unblocking Status: BANNED indicates that the ban has been lifted, RECOVER indicates that the ban has been lifted or not, BANNING indicates that the ban is in progress, RECOVERING indicates that the ban is being lifted, BAN_FAILED indicates that the ban has failed, RECOVER_FAILED indicates that the unblocking has failed.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"ip_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IP ListNote: This field may return null, indicating that a valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP.",
									},
									"provider": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Supplier, BGP represents default, CMCC represents China Mobile, CUCC represents China Unicom, and CTCC represents China Telecom.",
									},
									"bandwidth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Band width.",
									},
								},
							},
						},
						"http3_supported": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identification that supports the Http3 protocol, where:0 indicates shutdown;1 indicates enabled.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"in_ban_blacklist": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Is it on the banned blacklist? 0 indicates not on the blacklist, and 1 indicates on the blacklist.Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"feature_bitmap": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Property bitmap, where each bit represents a property, where:0 indicates that the feature is not supported;1, indicates support for this feature.The meaning of the feature bitmap is as follows (from right to left):The first bit supports 4-layer acceleration;The second bit supports 7-layer acceleration;The third bit supports Http3 access;The fourth bit supports IPv6;The fifth bit supports high-quality BGP access;The 6th bit supports three network access;The 7th bit supports QoS acceleration in the access segment.Note: This field may return null, indicating that a valid value cannot be obtained.Note: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudGaapProxyDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_proxy_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	proxyId := d.Get("proxy_id").(string)
	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var proxyDetail *gaap.ProxyInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapProxyDetail(ctx, proxyId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		proxyDetail = result
		return nil
	})
	if err != nil {
		return err
	}
	proxyInfoMap := map[string]interface{}{}

	if proxyDetail != nil {
		if proxyDetail.InstanceId != nil {
			proxyInfoMap["instance_id"] = proxyDetail.InstanceId
		}

		if proxyDetail.CreateTime != nil {
			proxyInfoMap["create_time"] = proxyDetail.CreateTime
		}

		if proxyDetail.ProjectId != nil {
			proxyInfoMap["project_id"] = proxyDetail.ProjectId
		}

		if proxyDetail.ProxyName != nil {
			proxyInfoMap["proxy_name"] = proxyDetail.ProxyName
		}

		if proxyDetail.AccessRegion != nil {
			proxyInfoMap["access_region"] = proxyDetail.AccessRegion
		}

		if proxyDetail.RealServerRegion != nil {
			proxyInfoMap["real_server_region"] = proxyDetail.RealServerRegion
		}

		if proxyDetail.Bandwidth != nil {
			proxyInfoMap["bandwidth"] = proxyDetail.Bandwidth
		}

		if proxyDetail.Concurrent != nil {
			proxyInfoMap["concurrent"] = proxyDetail.Concurrent
		}

		if proxyDetail.Status != nil {
			proxyInfoMap["status"] = proxyDetail.Status
		}

		if proxyDetail.Domain != nil {
			proxyInfoMap["domain"] = proxyDetail.Domain
		}

		if proxyDetail.IP != nil {
			proxyInfoMap["ip"] = proxyDetail.IP
		}

		if proxyDetail.Version != nil {
			proxyInfoMap["version"] = proxyDetail.Version
		}

		if proxyDetail.ProxyId != nil {
			proxyInfoMap["proxy_id"] = proxyDetail.ProxyId
		}

		if proxyDetail.Scalarable != nil {
			proxyInfoMap["scalarable"] = proxyDetail.Scalarable
		}

		if proxyDetail.SupportProtocols != nil {
			proxyInfoMap["support_protocols"] = proxyDetail.SupportProtocols
		}

		if proxyDetail.GroupId != nil {
			proxyInfoMap["group_id"] = proxyDetail.GroupId
		}

		if proxyDetail.PolicyId != nil {
			proxyInfoMap["policy_id"] = proxyDetail.PolicyId
		}

		if proxyDetail.AccessRegionInfo != nil {
			accessRegionInfoMap := map[string]interface{}{}

			if proxyDetail.AccessRegionInfo.RegionId != nil {
				accessRegionInfoMap["region_id"] = proxyDetail.AccessRegionInfo.RegionId
			}

			if proxyDetail.AccessRegionInfo.RegionName != nil {
				accessRegionInfoMap["region_name"] = proxyDetail.AccessRegionInfo.RegionName
			}

			if proxyDetail.AccessRegionInfo.RegionArea != nil {
				accessRegionInfoMap["region_area"] = proxyDetail.AccessRegionInfo.RegionArea
			}

			if proxyDetail.AccessRegionInfo.RegionAreaName != nil {
				accessRegionInfoMap["region_area_name"] = proxyDetail.AccessRegionInfo.RegionAreaName
			}

			if proxyDetail.AccessRegionInfo.IDCType != nil {
				accessRegionInfoMap["idc_type"] = proxyDetail.AccessRegionInfo.IDCType
			}

			if proxyDetail.AccessRegionInfo.FeatureBitmap != nil {
				accessRegionInfoMap["feature_bitmap"] = proxyDetail.AccessRegionInfo.FeatureBitmap
			}

			if proxyDetail.AccessRegionInfo.SupportFeature != nil {
				supportFeatureMap := map[string]interface{}{}

				if proxyDetail.AccessRegionInfo.SupportFeature.NetworkType != nil {
					supportFeatureMap["network_type"] = proxyDetail.AccessRegionInfo.SupportFeature.NetworkType
				}

				accessRegionInfoMap["support_feature"] = []interface{}{supportFeatureMap}
			}

			proxyInfoMap["access_region_info"] = []interface{}{accessRegionInfoMap}
		}

		if proxyDetail.RealServerRegionInfo != nil {
			realServerRegionInfoMap := map[string]interface{}{}

			if proxyDetail.RealServerRegionInfo.RegionId != nil {
				realServerRegionInfoMap["region_id"] = proxyDetail.RealServerRegionInfo.RegionId
			}

			if proxyDetail.RealServerRegionInfo.RegionName != nil {
				realServerRegionInfoMap["region_name"] = proxyDetail.RealServerRegionInfo.RegionName
			}

			if proxyDetail.RealServerRegionInfo.RegionArea != nil {
				realServerRegionInfoMap["region_area"] = proxyDetail.RealServerRegionInfo.RegionArea
			}

			if proxyDetail.RealServerRegionInfo.RegionAreaName != nil {
				realServerRegionInfoMap["region_area_name"] = proxyDetail.RealServerRegionInfo.RegionAreaName
			}

			if proxyDetail.RealServerRegionInfo.IDCType != nil {
				realServerRegionInfoMap["idc_type"] = proxyDetail.RealServerRegionInfo.IDCType
			}

			if proxyDetail.RealServerRegionInfo.FeatureBitmap != nil {
				realServerRegionInfoMap["feature_bitmap"] = proxyDetail.RealServerRegionInfo.FeatureBitmap
			}

			if proxyDetail.RealServerRegionInfo.SupportFeature != nil {
				supportFeatureMap := map[string]interface{}{}

				if proxyDetail.RealServerRegionInfo.SupportFeature.NetworkType != nil {
					supportFeatureMap["network_type"] = proxyDetail.RealServerRegionInfo.SupportFeature.NetworkType
				}

				realServerRegionInfoMap["support_feature"] = []interface{}{supportFeatureMap}
			}

			proxyInfoMap["real_server_region_info"] = []interface{}{realServerRegionInfoMap}
		}

		if proxyDetail.ForwardIP != nil {
			proxyInfoMap["forward_ip"] = proxyDetail.ForwardIP
		}

		if proxyDetail.TagSet != nil {
			tagSetList := []interface{}{}
			for _, tagSet := range proxyDetail.TagSet {
				tagSetMap := map[string]interface{}{}

				if tagSet.TagKey != nil {
					tagSetMap["tag_key"] = tagSet.TagKey
				}

				if tagSet.TagValue != nil {
					tagSetMap["tag_value"] = tagSet.TagValue
				}

				tagSetList = append(tagSetList, tagSetMap)
			}

			proxyInfoMap["tag_set"] = tagSetList
		}

		if proxyDetail.SupportSecurity != nil {
			proxyInfoMap["support_security"] = proxyDetail.SupportSecurity
		}

		if proxyDetail.BillingType != nil {
			proxyInfoMap["billing_type"] = proxyDetail.BillingType
		}

		if proxyDetail.RelatedGlobalDomains != nil {
			proxyInfoMap["related_global_domains"] = proxyDetail.RelatedGlobalDomains
		}

		if proxyDetail.ModifyConfigTime != nil {
			proxyInfoMap["modify_config_time"] = proxyDetail.ModifyConfigTime
		}

		if proxyDetail.ProxyType != nil {
			proxyInfoMap["proxy_type"] = proxyDetail.ProxyType
		}

		if proxyDetail.ClientIPMethod != nil {
			proxyInfoMap["client_ip_method"] = proxyDetail.ClientIPMethod
		}

		if proxyDetail.IPAddressVersion != nil {
			proxyInfoMap["ip_address_version"] = proxyDetail.IPAddressVersion
		}

		if proxyDetail.NetworkType != nil {
			proxyInfoMap["network_type"] = proxyDetail.NetworkType
		}

		if proxyDetail.PackageType != nil {
			proxyInfoMap["package_type"] = proxyDetail.PackageType
		}

		if proxyDetail.BanStatus != nil {
			proxyInfoMap["ban_status"] = proxyDetail.BanStatus
		}

		if proxyDetail.IPList != nil {
			iPListList := []interface{}{}
			for _, iPList := range proxyDetail.IPList {
				iPListMap := map[string]interface{}{}

				if iPList.IP != nil {
					iPListMap["ip"] = iPList.IP
				}

				if iPList.Provider != nil {
					iPListMap["provider"] = iPList.Provider
				}

				if iPList.Bandwidth != nil {
					iPListMap["bandwidth"] = iPList.Bandwidth
				}

				iPListList = append(iPListList, iPListMap)
			}

			proxyInfoMap["ip_list"] = iPListList
		}

		if proxyDetail.Http3Supported != nil {
			proxyInfoMap["http3_supported"] = proxyDetail.Http3Supported
		}

		if proxyDetail.InBanBlacklist != nil {
			proxyInfoMap["in_ban_blacklist"] = proxyDetail.InBanBlacklist
		}

		if proxyDetail.FeatureBitmap != nil {
			proxyInfoMap["feature_bitmap"] = proxyDetail.FeatureBitmap
		}

		_ = d.Set("proxy_detail", []interface{}{proxyInfoMap})
	}

	d.SetId(proxyId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), proxyInfoMap); e != nil {
			return e
		}
	}
	return nil
}

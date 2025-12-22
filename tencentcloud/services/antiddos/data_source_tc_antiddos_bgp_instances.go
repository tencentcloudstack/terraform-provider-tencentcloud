package antiddos

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddosv20200309 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAntiddosBgpInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAntiddosBgpInstancesRead,
		Schema: map[string]*schema.Schema{
			"filter_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region.",
			},

			"filter_instance_id_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Instance ID list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"filter_tag": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by tag key and value.",
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

			// computed
			"bgp_instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Returns purchased Anti-DDoS package information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"instance_charge_prepaid": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Renewal period related.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Purchase duration: unit in months.",
									},
									"renew_flag": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "NOTIFY_AND_MANUAL_RENEW: Notify expiration without automatic renewal.\nNOTIFY_AND_AUTO_RENEW: Notify expiration and automatically renew.\nDISABLE_NOTIFY_AND_MANUAL_RENEW: No notification and no automatic renewal.\nDefault: Notify expiration without automatic renewal.",
									},
								},
							},
						},
						"enterprise_package_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Enterprise edition Anti-DDoS package configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region where the Anti-DDoS package is purchased.",
									},
									"protect_ip_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of protected IPs.",
									},
									"basic_protect_bandwidth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Basic protection bandwidth.",
									},
									"bandwidth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Business bandwidth scale.",
									},
									"elastic_protect_bandwidth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Elastic bandwidth in Gbps, selectable elastic bandwidth [0,400,500,600,800,1000].\nDefault is 0.",
									},
									"elastic_bandwidth_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to enable elastic business bandwidth.\nDefault is false.",
									},
								},
							},
						},
						"standard_package_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Standard edition Anti-DDoS package configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region where the Anti-DDoS package is purchased.",
									},
									"protect_ip_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of protected IPs.",
									},
									"bandwidth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Protection business bandwidth 50Mbps.",
									},
									"elastic_bandwidth_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to enable elastic protection bandwidth. true: enable \nDefault is false: disable.",
									},
								},
							},
						},
						"standard_plus_package_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Standard edition 2.0 Anti-DDoS package configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region where the Anti-DDoS package is purchased.",
									},
									"protect_count": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protection count: TWO_TIMES: two full protections, UNLIMITED: unlimited protections.",
									},
									"protect_ip_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of protected IPs.",
									},
									"bandwidth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Protection bandwidth 50Mbps.",
									},
									"elastic_bandwidth_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to enable elastic business bandwidth.\ntrue: enable\nfalse: disable \nDefault is disable.",
									},
								},
							},
						},
						"tag_info_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag information.",
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
						"package_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Anti-DDoS package type.",
						},
						"instance_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Payment method.",
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

func dataSourceTencentCloudAntiddosBgpInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_antiddos_bgp_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(nil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service      = AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		filterRegion string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filter_region"); ok {
		paramMap["FilterRegion"] = helper.String(v.(string))
		filterRegion = v.(string)
	}

	if v, ok := d.GetOk("filter_instance_id_list"); ok {
		filterInstanceIdListList := []*string{}
		filterInstanceIdListSet := v.(*schema.Set).List()
		for i := range filterInstanceIdListSet {
			filterInstanceIdList := filterInstanceIdListSet[i].(string)
			filterInstanceIdListList = append(filterInstanceIdListList, helper.String(filterInstanceIdList))
		}

		paramMap["FilterInstanceIdList"] = filterInstanceIdListList
	}

	if v, ok := d.GetOk("filter_tag"); ok {
		filterTagSet := v.([]interface{})
		tmpSet := make([]*antiddosv20200309.TagInfo, 0, len(filterTagSet))
		for _, item := range filterTagSet {
			filterTagMap := item.(map[string]interface{})
			tagInfo := antiddosv20200309.TagInfo{}
			if v, ok := filterTagMap["tag_key"].(string); ok && v != "" {
				tagInfo.TagKey = helper.String(v)
			}

			if v, ok := filterTagMap["tag_value"].(string); ok && v != "" {
				tagInfo.TagValue = helper.String(v)
			}

			tmpSet = append(tmpSet, &tagInfo)
		}

		paramMap["FilterTag"] = tmpSet
	}

	var respData []*antiddosv20200309.BGPInstanceInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAntiddosBgpInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	bGPInstanceListList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, bGPInstanceList := range respData {
			bGPInstanceListMap := map[string]interface{}{}
			if bGPInstanceList.InstanceId != nil {
				bGPInstanceListMap["instance_id"] = bGPInstanceList.InstanceId
			}

			if bGPInstanceList.InstanceChargePrepaid != nil {
				instanceChargePrepaidMap := map[string]interface{}{}
				if bGPInstanceList.InstanceChargePrepaid.Period != nil {
					instanceChargePrepaidMap["period"] = bGPInstanceList.InstanceChargePrepaid.Period
				}

				if bGPInstanceList.InstanceChargePrepaid.RenewFlag != nil {
					instanceChargePrepaidMap["renew_flag"] = bGPInstanceList.InstanceChargePrepaid.RenewFlag
				}

				bGPInstanceListMap["instance_charge_prepaid"] = []interface{}{instanceChargePrepaidMap}
			}

			if bGPInstanceList.EnterprisePackageConfig != nil {
				enterprisePackageConfigMap := map[string]interface{}{}
				if bGPInstanceList.EnterprisePackageConfig.Region != nil {
					enterprisePackageConfigMap["region"] = bGPInstanceList.EnterprisePackageConfig.Region
				}

				if bGPInstanceList.EnterprisePackageConfig.ProtectIpCount != nil {
					enterprisePackageConfigMap["protect_ip_count"] = bGPInstanceList.EnterprisePackageConfig.ProtectIpCount
				}

				if bGPInstanceList.EnterprisePackageConfig.BasicProtectBandwidth != nil {
					enterprisePackageConfigMap["basic_protect_bandwidth"] = bGPInstanceList.EnterprisePackageConfig.BasicProtectBandwidth
				}

				if bGPInstanceList.EnterprisePackageConfig.Bandwidth != nil {
					enterprisePackageConfigMap["bandwidth"] = bGPInstanceList.EnterprisePackageConfig.Bandwidth
				}

				if bGPInstanceList.EnterprisePackageConfig.ElasticProtectBandwidth != nil {
					enterprisePackageConfigMap["elastic_protect_bandwidth"] = bGPInstanceList.EnterprisePackageConfig.ElasticProtectBandwidth
				}

				if bGPInstanceList.EnterprisePackageConfig.ElasticBandwidthFlag != nil {
					enterprisePackageConfigMap["elastic_bandwidth_flag"] = bGPInstanceList.EnterprisePackageConfig.ElasticBandwidthFlag
				}

				bGPInstanceListMap["enterprise_package_config"] = []interface{}{enterprisePackageConfigMap}
			}

			if bGPInstanceList.StandardPackageConfig != nil {
				standardPackageConfigMap := map[string]interface{}{}
				if bGPInstanceList.StandardPackageConfig.Region != nil {
					standardPackageConfigMap["region"] = bGPInstanceList.StandardPackageConfig.Region
				}

				if bGPInstanceList.StandardPackageConfig.ProtectIpCount != nil {
					standardPackageConfigMap["protect_ip_count"] = bGPInstanceList.StandardPackageConfig.ProtectIpCount
				}

				if bGPInstanceList.StandardPackageConfig.Bandwidth != nil {
					standardPackageConfigMap["bandwidth"] = bGPInstanceList.StandardPackageConfig.Bandwidth
				}

				if bGPInstanceList.StandardPackageConfig.ElasticBandwidthFlag != nil {
					standardPackageConfigMap["elastic_bandwidth_flag"] = bGPInstanceList.StandardPackageConfig.ElasticBandwidthFlag
				}

				bGPInstanceListMap["standard_package_config"] = []interface{}{standardPackageConfigMap}
			}

			if bGPInstanceList.StandardPlusPackageConfig != nil {
				standardPlusPackageConfigMap := map[string]interface{}{}
				if bGPInstanceList.StandardPlusPackageConfig.Region != nil {
					standardPlusPackageConfigMap["region"] = bGPInstanceList.StandardPlusPackageConfig.Region
				}

				if bGPInstanceList.StandardPlusPackageConfig.ProtectCount != nil {
					standardPlusPackageConfigMap["protect_count"] = bGPInstanceList.StandardPlusPackageConfig.ProtectCount
				}

				if bGPInstanceList.StandardPlusPackageConfig.ProtectIpCount != nil {
					standardPlusPackageConfigMap["protect_ip_count"] = bGPInstanceList.StandardPlusPackageConfig.ProtectIpCount
				}

				if bGPInstanceList.StandardPlusPackageConfig.Bandwidth != nil {
					standardPlusPackageConfigMap["bandwidth"] = bGPInstanceList.StandardPlusPackageConfig.Bandwidth
				}

				if bGPInstanceList.StandardPlusPackageConfig.ElasticBandwidthFlag != nil {
					standardPlusPackageConfigMap["elastic_bandwidth_flag"] = bGPInstanceList.StandardPlusPackageConfig.ElasticBandwidthFlag
				}

				bGPInstanceListMap["standard_plus_package_config"] = []interface{}{standardPlusPackageConfigMap}
			}

			tagInfoListList := make([]map[string]interface{}, 0, len(bGPInstanceList.TagInfoList))
			if bGPInstanceList.TagInfoList != nil {
				for _, tagInfoList := range bGPInstanceList.TagInfoList {
					tagInfoListMap := map[string]interface{}{}
					if tagInfoList.TagKey != nil {
						tagInfoListMap["tag_key"] = tagInfoList.TagKey
					}

					if tagInfoList.TagValue != nil {
						tagInfoListMap["tag_value"] = tagInfoList.TagValue
					}

					tagInfoListList = append(tagInfoListList, tagInfoListMap)
				}

				bGPInstanceListMap["tag_info_list"] = tagInfoListList
			}

			if bGPInstanceList.PackageType != nil {
				bGPInstanceListMap["package_type"] = bGPInstanceList.PackageType
			}

			if bGPInstanceList.InstanceChargeType != nil {
				bGPInstanceListMap["instance_charge_type"] = bGPInstanceList.InstanceChargeType
			}

			bGPInstanceListList = append(bGPInstanceListList, bGPInstanceListMap)
		}

		_ = d.Set("bgp_instance_list", bGPInstanceListList)
	}

	d.SetId(filterRegion)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

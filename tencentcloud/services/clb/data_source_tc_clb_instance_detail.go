package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbInstanceDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbInstanceDetailRead,
		Schema: map[string]*schema.Schema{
			"fields": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of fields. Only fields specified will be returned. If it's left blank, `null` is returned. The fields `LoadBalancerId` and `LoadBalancerName` are added by default. For details about fields.",
			},

			"target_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Target type. Valid values: NODE and GROUP. If the list of fields contains `TargetId`, `TargetAddress`, `TargetPort`, `TargetWeight` and other fields, `Target` of the target group or non-target group must be exported.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter condition of querying lists describing CLB instance details:loadbalancer-id - String - Required: no - (Filter condition) CLB instance ID, such as lb-12345678; project-id - String - Required: no - (Filter condition) Project ID, such as 0 and 123; network - String - Required: no - (Filter condition) Network type of the CLB instance, such as Public and Private.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt; vip - String - Required: no - (Filter condition) CLB instance VIP, such as 1.1.1.1 and 2204::22:3; target-ip - String - Required: no - (Filter condition) Private IP of the target real servers, such as1.1.1.1 and 2203::214:4; vpcid - String - Required: no - (Filter condition) Identifier of the VPC instance to which the CLB instance belongs, such as vpc-12345678; zone - String - Required: no - (Filter condition) Availability zone where the CLB instance resides, such as ap-guangzhou-1; tag-key - String - Required: no - (Filter condition) Tag key of the CLB instance, such as name; tag:* - String - Required: no - (Filter condition) CLB instance tag, followed by tag key after the colon. For example, use {Name: tag:name,Values: [zhangsan, lisi]} to filter the tag key `name` with the tag value `zhangsan` and `lisi`; fuzzy-search - String - Required: no - (Filter condition) Fuzzy search for CLB instance VIP and CLB instance name, such as 1.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter value array.",
						},
					},
				},
			},

			"load_balancer_detail_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of CLB instance details.Note: this field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance ID.",
						},
						"load_balancer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance name.",
						},
						"load_balancer_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance network type:Public: public network; Private: private network.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB instance status, including:0: creating; 1: running.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance VIP.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"address_ipv6": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv6 VIP address of the CLB instance.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"address_ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version of the CLB instance. Valid values: IPv4, IPv6.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"ipv6_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv6 address type of the CLB instance. Valid values: IPv6Nat64, IPv6FullChain.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone where the CLB instance resides.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"address_isp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ISP to which the CLB IP address belongs.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC instance to which the CLB instance belongs.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the project to which the CLB instance belongs. 0: default project.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance creation time.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB instance billing mode.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"network_attributes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CLB instance network attribute.Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"internet_charge_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TRAFFIC_POSTPAID_BY_HOUR: hourly pay-as-you-go by traffic; BANDWIDTH_POSTPAID_BY_HOUR: hourly pay-as-you-go by bandwidth;BANDWIDTH_PACKAGE: billed by bandwidth package (currently, this method is supported only if the ISP is specified).",
									},
									"internet_max_bandwidth_out": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum outbound bandwidth in Mbps, which applies only to public network CLB. Value range: 0-65,535. Default value: 10.",
									},
									"bandwidth_pkg_sub_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Bandwidth package type, such as SINGLEISPNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"prepaid_attributes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Pay-as-you-go attribute of the CLB instance.Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"renew_flag": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Renewal type. AUTO_RENEW: automatic renewal; MANUAL_RENEW: manual renewalNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cycle, indicating the number of months (reserved field)Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"extra_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Reserved field, which can be ignored generally.Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zhi_tong": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to enable VIP direct connectionNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"tgw_group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TgwGroup nameNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"config_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom configuration IDs of CLB instances. Multiple IDs must be separated by commas (,).Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CLB instance tag information.Note: this field may return null, indicating that no valid values can be obtained.",
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
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB listener ID.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Listener protocol.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Listener port.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"location_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forwarding rule ID.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name of the forwarding rule.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forwarding rule path.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"target_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of target real servers.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"target_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Address of target real servers.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"target_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Listening port of target real servers.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"target_weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Forwarding weight of target real servers.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"isolation": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "0: not isolated; 1: isolated.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"security_group": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of the security groups bound to the CLB instance.Note: this field may return `null`, indicating that no valid values can be obtained.",
						},
						"load_balancer_pass_to_target": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the CLB instance is billed by IP.Note: this field may return `null`, indicating that no valid values can be obtained.",
						},
						"target_health": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health status of the target real server.Note: this field may return `null`, indicating that no valid values can be obtained.",
						},
						"domains": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "List o domain names associated with the forwarding ruleNote: This field may return `null`, indicating that no valid values can be obtained.",
						},
						"slave_zone": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The secondary zone of multi-AZ CLB instanceNote: This field may return `null`, indicating that no valid values can be obtained.",
						},
						"zones": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The AZ of private CLB instance. This is only available for beta users.Note: This field may return `null`, indicating that no valid values can be obtained.",
						},
						"sni_switch": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether SNI is enabled. This parameter is only meaningful for HTTPS listeners.Note: This field may return `null`, indicating that no valid values can be obtained.",
						},
						"load_balancer_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name of the CLB instance.Note: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudClbInstanceDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_instance_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("fields"); ok {
		fieldsSet := v.(*schema.Set).List()
		paramMap["Fields"] = helper.InterfacesStringsPoint(fieldsSet)
	}

	if v, ok := d.GetOk("target_type"); ok {
		paramMap["TargetType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*clb.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := clb.Filter{}
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
		paramMap["Filters"] = tmpSet
	}

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var loadBalancerDetailSet []*clb.LoadBalancerDetail

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbInstanceDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		loadBalancerDetailSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(loadBalancerDetailSet))
	tmpList := make([]map[string]interface{}, 0, len(loadBalancerDetailSet))

	if loadBalancerDetailSet != nil {
		for _, loadBalancerDetail := range loadBalancerDetailSet {
			loadBalancerDetailMap := map[string]interface{}{}

			if loadBalancerDetail.LoadBalancerId != nil {
				loadBalancerDetailMap["load_balancer_id"] = loadBalancerDetail.LoadBalancerId
			}

			if loadBalancerDetail.LoadBalancerName != nil {
				loadBalancerDetailMap["load_balancer_name"] = loadBalancerDetail.LoadBalancerName
			}

			if loadBalancerDetail.LoadBalancerType != nil {
				loadBalancerDetailMap["load_balancer_type"] = loadBalancerDetail.LoadBalancerType
			}

			if loadBalancerDetail.Status != nil {
				loadBalancerDetailMap["status"] = loadBalancerDetail.Status
			}

			if loadBalancerDetail.Address != nil {
				loadBalancerDetailMap["address"] = loadBalancerDetail.Address
			}

			if loadBalancerDetail.AddressIPv6 != nil {
				loadBalancerDetailMap["address_ipv6"] = loadBalancerDetail.AddressIPv6
			}

			if loadBalancerDetail.AddressIPVersion != nil {
				loadBalancerDetailMap["address_ip_version"] = loadBalancerDetail.AddressIPVersion
			}

			if loadBalancerDetail.IPv6Mode != nil {
				loadBalancerDetailMap["ipv6_mode"] = loadBalancerDetail.IPv6Mode
			}

			if loadBalancerDetail.Zone != nil {
				loadBalancerDetailMap["zone"] = loadBalancerDetail.Zone
			}

			if loadBalancerDetail.AddressIsp != nil {
				loadBalancerDetailMap["address_isp"] = loadBalancerDetail.AddressIsp
			}

			if loadBalancerDetail.VpcId != nil {
				loadBalancerDetailMap["vpc_id"] = loadBalancerDetail.VpcId
			}

			if loadBalancerDetail.ProjectId != nil {
				loadBalancerDetailMap["project_id"] = loadBalancerDetail.ProjectId
			}

			if loadBalancerDetail.CreateTime != nil {
				loadBalancerDetailMap["create_time"] = loadBalancerDetail.CreateTime
			}

			if loadBalancerDetail.ChargeType != nil {
				loadBalancerDetailMap["charge_type"] = loadBalancerDetail.ChargeType
			}

			if loadBalancerDetail.NetworkAttributes != nil {
				networkAttributesMap := map[string]interface{}{}

				if loadBalancerDetail.NetworkAttributes.InternetChargeType != nil {
					networkAttributesMap["internet_charge_type"] = loadBalancerDetail.NetworkAttributes.InternetChargeType
				}

				if loadBalancerDetail.NetworkAttributes.InternetMaxBandwidthOut != nil {
					networkAttributesMap["internet_max_bandwidth_out"] = loadBalancerDetail.NetworkAttributes.InternetMaxBandwidthOut
				}

				if loadBalancerDetail.NetworkAttributes.BandwidthpkgSubType != nil {
					networkAttributesMap["bandwidth_pkg_sub_type"] = loadBalancerDetail.NetworkAttributes.BandwidthpkgSubType
				}

				loadBalancerDetailMap["network_attributes"] = []interface{}{networkAttributesMap}
			}

			if loadBalancerDetail.PrepaidAttributes != nil {
				prepaidAttributesMap := map[string]interface{}{}

				if loadBalancerDetail.PrepaidAttributes.RenewFlag != nil {
					prepaidAttributesMap["renew_flag"] = loadBalancerDetail.PrepaidAttributes.RenewFlag
				}

				if loadBalancerDetail.PrepaidAttributes.Period != nil {
					prepaidAttributesMap["period"] = loadBalancerDetail.PrepaidAttributes.Period
				}

				loadBalancerDetailMap["prepaid_attributes"] = []interface{}{prepaidAttributesMap}
			}

			if loadBalancerDetail.ExtraInfo != nil {
				extraInfoMap := map[string]interface{}{}

				if loadBalancerDetail.ExtraInfo.ZhiTong != nil {
					extraInfoMap["zhi_tong"] = loadBalancerDetail.ExtraInfo.ZhiTong
				}

				if loadBalancerDetail.ExtraInfo.TgwGroupName != nil {
					extraInfoMap["tgw_group_name"] = loadBalancerDetail.ExtraInfo.TgwGroupName
				}

				loadBalancerDetailMap["extra_info"] = []interface{}{extraInfoMap}
			}

			if loadBalancerDetail.ConfigId != nil {
				loadBalancerDetailMap["config_id"] = loadBalancerDetail.ConfigId
			}

			if loadBalancerDetail.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range loadBalancerDetail.Tags {
					tagsMap := map[string]interface{}{}

					if tags.TagKey != nil {
						tagsMap["tag_key"] = tags.TagKey
					}

					if tags.TagValue != nil {
						tagsMap["tag_value"] = tags.TagValue
					}

					tagsList = append(tagsList, tagsMap)
				}

				loadBalancerDetailMap["tags"] = tagsList
			}

			if loadBalancerDetail.ListenerId != nil {
				loadBalancerDetailMap["listener_id"] = loadBalancerDetail.ListenerId
			}

			if loadBalancerDetail.Protocol != nil {
				loadBalancerDetailMap["protocol"] = loadBalancerDetail.Protocol
			}

			if loadBalancerDetail.Port != nil {
				loadBalancerDetailMap["port"] = loadBalancerDetail.Port
			}

			if loadBalancerDetail.LocationId != nil {
				loadBalancerDetailMap["location_id"] = loadBalancerDetail.LocationId
			}

			if loadBalancerDetail.Domain != nil {
				loadBalancerDetailMap["domain"] = loadBalancerDetail.Domain
			}

			if loadBalancerDetail.Url != nil {
				loadBalancerDetailMap["url"] = loadBalancerDetail.Url
			}

			if loadBalancerDetail.TargetId != nil {
				loadBalancerDetailMap["target_id"] = loadBalancerDetail.TargetId
			}

			if loadBalancerDetail.TargetAddress != nil {
				loadBalancerDetailMap["target_address"] = loadBalancerDetail.TargetAddress
			}

			if loadBalancerDetail.TargetPort != nil {
				loadBalancerDetailMap["target_port"] = loadBalancerDetail.TargetPort
			}

			if loadBalancerDetail.TargetWeight != nil {
				loadBalancerDetailMap["target_weight"] = loadBalancerDetail.TargetWeight
			}

			if loadBalancerDetail.Isolation != nil {
				loadBalancerDetailMap["isolation"] = loadBalancerDetail.Isolation
			}

			if loadBalancerDetail.SecurityGroup != nil {
				loadBalancerDetailMap["security_group"] = loadBalancerDetail.SecurityGroup
			}

			if loadBalancerDetail.LoadBalancerPassToTarget != nil {
				loadBalancerDetailMap["load_balancer_pass_to_target"] = loadBalancerDetail.LoadBalancerPassToTarget
			}

			if loadBalancerDetail.TargetHealth != nil {
				loadBalancerDetailMap["target_health"] = loadBalancerDetail.TargetHealth
			}

			if loadBalancerDetail.Domains != nil {
				loadBalancerDetailMap["domains"] = loadBalancerDetail.Domains
			}

			if loadBalancerDetail.SlaveZone != nil {
				loadBalancerDetailMap["slave_zone"] = loadBalancerDetail.SlaveZone
			}

			if loadBalancerDetail.Zones != nil {
				loadBalancerDetailMap["zones"] = loadBalancerDetail.Zones
			}

			if loadBalancerDetail.SniSwitch != nil {
				loadBalancerDetailMap["sni_switch"] = loadBalancerDetail.SniSwitch
			}

			if loadBalancerDetail.LoadBalancerDomain != nil {
				loadBalancerDetailMap["load_balancer_domain"] = loadBalancerDetail.LoadBalancerDomain
			}

			ids = append(ids, *loadBalancerDetail.LoadBalancerId)
			tmpList = append(tmpList, loadBalancerDetailMap)
		}

		_ = d.Set("load_balancer_detail_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}

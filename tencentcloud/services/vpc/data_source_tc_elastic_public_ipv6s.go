package vpc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudElasticPublicIpv6s() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudElasticPublicIpv6sRead,
		Schema: map[string]*schema.Schema{
			"ipv6_address_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Unique ID column that identifies IPv6.\n\t- Traditional Elastic IPv6 unique ID is like: `eip-11112222`\n\t- Elastic IPv6 unique ID is like: `eipv6 -11112222`\nNote: Parameters do not support specifying both IPv6AddressIds and Filters.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The detailed filter conditions are as follows:\n\t- address-id-String-required: no-(filter condition) filter by the unique ID of the elastic public network IPv6.\n\t- public-ipv6-address-String-required: no-(filter condition) filter by the IP address of the public network IPv6.\n\t- charge-type-String-required: no-(filter condition) filter by billing type.\n\t- private-ipv6-address-String-required: no-(filter condition) filter by bound private network IPv6 address.\n\t- egress-String-required: no-(filter condition) filter by exit.\n\t- address-type-String-required: no-(filter condition) filter by IPv6 type.\n\t- address-isp-String-required: no-(filter condition) filter by operator type.\n  The status includes: 'CREATING','BINDING','BIND','UNBINDING','UNBIND','OFFLINING','BIND_ENI','PRIVATE'.\n\t- address-name-String-required: no-(filter condition) filter by EIP name. Blur filtering is not supported.\n\t- tag-key-String-required: no-(filter condition) filter by label key.\n\t- tag-value-String-required: no-(filter condition) filter by tag value.\n\t- tag:tag-key-String-required: no-(filter condition) filter by label key value pair. Tag-key is replaced with a specific label key.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Property name. If there are multiple Filters, the relationship between Filters is a logical AND (AND) relationship.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Attribute value. If there are multiple Values in the same Filter, the relationship between Values under the same Filter is a logical OR relationship. When the value type is a Boolean type, the value can be directly taken to the string TRUE or FALSE.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"traditional": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to query traditional IPv6 address information.",
			},

			"address_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of IPv6 details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of EIP is the unique identifier of EIP.",
						},
						"address_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "EIP name.",
						},
						"address_status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "EIP status, including CREATING(Creating),BINDING(Binding),BIND(Binding),UNBINDING(Unbinding),UNBIND(Unbinding),OFFLINING(Releasing),BIND_ENI(Binding Suspend Elastic Network Interface).",
						},
						"address_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "External network IP address.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The bound resource instance `ID`. It may be a `CVM`,`NAT`.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Creation time. It is expressed in accordance with the ISO8601 standard and uses UTC time. The format is: `Y-MM-DDThh:mm:ssZ`.",
						},
						"network_interface_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The bound Elastic Network Interface ID.",
						},
						"private_address_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Binding resources intranet IP.",
						},
						"is_arrears": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Resource isolation status. true means eip is in isolation, false means resource is in non-isolation state.",
						},
						"is_blocked": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Resource blocking status. true means eip is blocked, false means eip is not blocked.",
						},
						"is_eip_direct_connection": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether eip supports pass-through mode. true means eip supports pass-through mode, false means resources do not support pass-through mode.",
						},
						"address_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "EIP resource types, including CalcIP, WanIP, EIP, AnycastEIP, and high-defense EIP. Among them: `CalcIP` means device IP,`WanIP` means ordinary public network IP,`EIP` means elastic public network IP,`AnycastEIP` means accelerated EIP, and `AntiDDoSEIP` means highly resistant EIP.",
						},
						"cascade_release": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether eip is automatically released after unbinding. true means that eip will be automatically released after unbinding, false means that eip will not be automatically released after unbinding.",
						},
						"eip_alg_type": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "The type of protocol opened by EIP ALG.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ftp": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether the Ftp protocol Alg function is enabled.",
									},
									"sip": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether the Sip protocol Alg function is enabled.",
									},
								},
							},
						},
						"internet_service_provider": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Operator information of elastic public network IP. Current possible return values include `CMCC`,`CTCC`,`CUCC`,`BGP`.",
						},
						"local_bgp": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether local bandwidth EIP.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The bandwidth value of the elastic public network IP. Note that the elastic public IP of traditional account types has no bandwidth attribute and the value is null.",
						},
						"internet_charge_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Network charging model for elastic public network IP. Note that the elastic public IP of traditional account types does not have a network charging mode attribute and the value is blank. Note: This field may return null, indicating that a valid value cannot be obtained. Includes: \nBANDWIDTH_PREPAID_BY_MONTH: indicates a prepaid monthly bandwidth. \nTRAFFIC_POSTPAID_BY_HOUR: means post-payment per hour. BANDWIDTH_POSTPAID_BY_HOUR: means postpayment per hour of bandwidth.\nBANDWIDTH_PACKAGE: indicates a shared Bandwidth Package.",
						},
						"tag_set": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of tags associated with elastic public IP.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag key",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Tag value",
									},
								},
							},
						},
						"deadline_date": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Expiration time.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The instance type of the EIP binding.",
						},
						"egress": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Static single-wire IP network exit.",
						},
						"anti_ddos_package_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "High-defense package ID. When the EIP type is high-defense EIP, it returns the high-defense package ID to which the EIP is bound.",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether the current EIP is automatically renewed, this field will be displayed only for EIP prepaid by monthly bandwidth. Examples of specific values are as follows:\n\t- NOTIFY_AND_MANUAL_RENEW: Normal renewal\n\t- NOTIFY_AND_AUTO_RENEW: Automatic renewal\n\t- DISABLE_NOTIFY_AND_MANUAL_RENEW: No renewal after expiration.",
						},
						"bandwidth_package_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The bandwidth package ID associated with the current public IP. If the public IP does not use bandwidth packages for charging, the return will be blank.",
						},
						"un_vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique ID of the vpc to which traditional Elastic IPv6 belongs.",
						},
						"dedicated_cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "CDC unique ID.",
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

func dataSourceTencentCloudElasticPublicIpv6sRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_elastic_public_ipv6s.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ipv6_address_ids"); ok {
		iPv6AddressIdsList := []*string{}
		iPv6AddressIdsSet := v.(*schema.Set).List()
		for i := range iPv6AddressIdsSet {
			iPv6AddressIds := iPv6AddressIdsSet[i].(string)
			iPv6AddressIdsList = append(iPv6AddressIdsList, helper.String(iPv6AddressIds))
		}
		paramMap["IPv6AddressIds"] = iPv6AddressIdsList
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*vpc.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := vpc.Filter{}
			if v, ok := filtersMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(values))
				}
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOkExists("traditional"); ok {
		paramMap["Traditional"] = helper.Bool(v.(bool))
	}

	var respData []*vpc.Address
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeElasticPublicIpv6sByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	var ids []string
	addressSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, addressSet := range respData {
			addressSetMap := map[string]interface{}{}

			var addressId string
			if addressSet.AddressId != nil {
				addressSetMap["address_id"] = addressSet.AddressId
				addressId = *addressSet.AddressId
			}

			if addressSet.AddressName != nil {
				addressSetMap["address_name"] = addressSet.AddressName
			}

			if addressSet.AddressStatus != nil {
				addressSetMap["address_status"] = addressSet.AddressStatus
			}

			if addressSet.AddressIp != nil {
				addressSetMap["address_ip"] = addressSet.AddressIp
			}

			if addressSet.InstanceId != nil {
				addressSetMap["instance_id"] = addressSet.InstanceId
			}

			if addressSet.CreatedTime != nil {
				addressSetMap["created_time"] = addressSet.CreatedTime
			}

			if addressSet.NetworkInterfaceId != nil {
				addressSetMap["network_interface_id"] = addressSet.NetworkInterfaceId
			}

			if addressSet.PrivateAddressIp != nil {
				addressSetMap["private_address_ip"] = addressSet.PrivateAddressIp
			}

			if addressSet.IsArrears != nil {
				addressSetMap["is_arrears"] = addressSet.IsArrears
			}

			if addressSet.IsBlocked != nil {
				addressSetMap["is_blocked"] = addressSet.IsBlocked
			}

			if addressSet.IsEipDirectConnection != nil {
				addressSetMap["is_eip_direct_connection"] = addressSet.IsEipDirectConnection
			}

			if addressSet.AddressType != nil {
				addressSetMap["address_type"] = addressSet.AddressType
			}

			if addressSet.CascadeRelease != nil {
				addressSetMap["cascade_release"] = addressSet.CascadeRelease
			}

			eipAlgTypeMap := map[string]interface{}{}

			if addressSet.EipAlgType != nil {
				if addressSet.EipAlgType.Ftp != nil {
					eipAlgTypeMap["ftp"] = addressSet.EipAlgType.Ftp
				}

				if addressSet.EipAlgType.Sip != nil {
					eipAlgTypeMap["sip"] = addressSet.EipAlgType.Sip
				}

				addressSetMap["eip_alg_type"] = []interface{}{eipAlgTypeMap}
			}

			if addressSet.InternetServiceProvider != nil {
				addressSetMap["internet_service_provider"] = addressSet.InternetServiceProvider
			}

			if addressSet.LocalBgp != nil {
				addressSetMap["local_bgp"] = addressSet.LocalBgp
			}

			if addressSet.Bandwidth != nil {
				addressSetMap["bandwidth"] = addressSet.Bandwidth
			}

			if addressSet.InternetChargeType != nil {
				addressSetMap["internet_charge_type"] = addressSet.InternetChargeType
			}

			tagSetList := make([]map[string]interface{}, 0, len(addressSet.TagSet))
			if addressSet.TagSet != nil {
				for _, tagSet := range addressSet.TagSet {
					tagSetMap := map[string]interface{}{}

					if tagSet.Key != nil {
						tagSetMap["key"] = tagSet.Key
					}

					if tagSet.Value != nil {
						tagSetMap["value"] = tagSet.Value
					}

					tagSetList = append(tagSetList, tagSetMap)
				}

				addressSetMap["tag_set"] = tagSetList
			}
			if addressSet.DeadlineDate != nil {
				addressSetMap["deadline_date"] = addressSet.DeadlineDate
			}

			if addressSet.InstanceType != nil {
				addressSetMap["instance_type"] = addressSet.InstanceType
			}

			if addressSet.Egress != nil {
				addressSetMap["egress"] = addressSet.Egress
			}

			if addressSet.AntiDDoSPackageId != nil {
				addressSetMap["anti_ddos_package_id"] = addressSet.AntiDDoSPackageId
			}

			if addressSet.RenewFlag != nil {
				addressSetMap["renew_flag"] = addressSet.RenewFlag
			}

			if addressSet.BandwidthPackageId != nil {
				addressSetMap["bandwidth_package_id"] = addressSet.BandwidthPackageId
			}

			if addressSet.UnVpcId != nil {
				addressSetMap["un_vpc_id"] = addressSet.UnVpcId
			}

			if addressSet.DedicatedClusterId != nil {
				addressSetMap["dedicated_cluster_id"] = addressSet.DedicatedClusterId
			}

			ids = append(ids, addressId)
			addressSetList = append(addressSetList, addressSetMap)
		}

		_ = d.Set("address_set", addressSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), addressSetList); e != nil {
			return e
		}
	}

	return nil
}

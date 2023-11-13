/*
Use this data source to query detailed information of tse gateway_config

Example Usage

```hcl
data "tencentcloud_tse_gateway_config" "gateway_config" {
  gateway_id = "gateway-xx"
  group_id = "group-xx"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTseGatewayConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGatewayConfigRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway ID.",
			},

			"group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Gateway group ID, default value: default group ID.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Gateway network informations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway ID.",
						},
						"config_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Network informations of group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"console_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of console.",
									},
									"http_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Http urlNote: This field may return null, indicating that a valid value is not available.",
									},
									"https_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Https url.",
									},
									"net_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network type. Reference value:- Open- Internal.",
									},
									"admin_user": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Administrator user name of kongNote: This field may return null, indicating that a valid value is not available.",
									},
									"admin_password": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Administrator password of kongNote: This field may return null, indicating that a valid value is not available.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network status. Reference value:- Open- Closed- UpdatingNote: This field may return null, indicating that a valid value is not available.",
									},
									"access_control": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Network access policyNote: This field may return null, indicating that a valid value is not available.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Access type. Reference value:- Whitelist- Blacklist.",
												},
												"cidr_white_list": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Whitelist.",
												},
												"cidr_black_list": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Blacklist.",
												},
											},
										},
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet IDNote: This field may return null, indicating that a valid value is not available.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Vpc IDNote: This field may return null, indicating that a valid value is not available.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DescriptionNote: This field may return null, indicating that a valid value is not available.",
									},
									"sla_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specification type of clb, return SLA means performance capacity type，return empty means shared typeNote: This field may return null, indicating that a valid value is not available.",
									},
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Clb vipNote: This field may return null, indicating that a valid value is not available.",
									},
									"internet_max_bandwidth_out": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Public network outbound traffic bandwidthNote: This field may return null, indicating that a valid value is not available.",
									},
									"multi_zone_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether load balancing has multiple availability zonesNote: This field may return null, indicating that a valid value is not available.",
									},
									"master_zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Primary availability zoneNote: This field may return null, indicating that a valid value is not available.",
									},
									"slave_zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Alternate availability zoneNote: This field may return null, indicating that a valid value is not available.",
									},
									"master_zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of primary availability zoneNote: This field may return null, indicating that a valid value is not available.",
									},
									"slave_zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of alternate availability zoneNote: This field may return null, indicating that a valid value is not available.",
									},
									"network_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network IDNote: This field may return null, indicating that a valid value is not available.",
									},
								},
							},
						},
						"group_subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID of groupNote: This field may return null, indicating that a valid value is not available.",
						},
						"group_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vpc ID of groupNote: This field may return null, indicating that a valid value is not available.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group IDNote: This field may return null, indicating that a valid value is not available.",
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

func dataSourceTencentCloudTseGatewayConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tse_gateway_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*tse.DescribeCloudNativeAPIGatewayConfigResult

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTseGatewayConfigByFilter(ctx, paramMap)
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
		describeCloudNativeAPIGatewayConfigResultMap := map[string]interface{}{}

		if result.GatewayId != nil {
			describeCloudNativeAPIGatewayConfigResultMap["gateway_id"] = result.GatewayId
		}

		if result.ConfigList != nil {
			configListList := []interface{}{}
			for _, configList := range result.ConfigList {
				configListMap := map[string]interface{}{}

				if configList.ConsoleType != nil {
					configListMap["console_type"] = configList.ConsoleType
				}

				if configList.HttpUrl != nil {
					configListMap["http_url"] = configList.HttpUrl
				}

				if configList.HttpsUrl != nil {
					configListMap["https_url"] = configList.HttpsUrl
				}

				if configList.NetType != nil {
					configListMap["net_type"] = configList.NetType
				}

				if configList.AdminUser != nil {
					configListMap["admin_user"] = configList.AdminUser
				}

				if configList.AdminPassword != nil {
					configListMap["admin_password"] = configList.AdminPassword
				}

				if configList.Status != nil {
					configListMap["status"] = configList.Status
				}

				if configList.AccessControl != nil {
					accessControlMap := map[string]interface{}{}

					if configList.AccessControl.Mode != nil {
						accessControlMap["mode"] = configList.AccessControl.Mode
					}

					if configList.AccessControl.CidrWhiteList != nil {
						accessControlMap["cidr_white_list"] = configList.AccessControl.CidrWhiteList
					}

					if configList.AccessControl.CidrBlackList != nil {
						accessControlMap["cidr_black_list"] = configList.AccessControl.CidrBlackList
					}

					configListMap["access_control"] = []interface{}{accessControlMap}
				}

				if configList.SubnetId != nil {
					configListMap["subnet_id"] = configList.SubnetId
				}

				if configList.VpcId != nil {
					configListMap["vpc_id"] = configList.VpcId
				}

				if configList.Description != nil {
					configListMap["description"] = configList.Description
				}

				if configList.SlaType != nil {
					configListMap["sla_type"] = configList.SlaType
				}

				if configList.Vip != nil {
					configListMap["vip"] = configList.Vip
				}

				if configList.InternetMaxBandwidthOut != nil {
					configListMap["internet_max_bandwidth_out"] = configList.InternetMaxBandwidthOut
				}

				if configList.MultiZoneFlag != nil {
					configListMap["multi_zone_flag"] = configList.MultiZoneFlag
				}

				if configList.MasterZoneId != nil {
					configListMap["master_zone_id"] = configList.MasterZoneId
				}

				if configList.SlaveZoneId != nil {
					configListMap["slave_zone_id"] = configList.SlaveZoneId
				}

				if configList.MasterZoneName != nil {
					configListMap["master_zone_name"] = configList.MasterZoneName
				}

				if configList.SlaveZoneName != nil {
					configListMap["slave_zone_name"] = configList.SlaveZoneName
				}

				if configList.NetworkId != nil {
					configListMap["network_id"] = configList.NetworkId
				}

				configListList = append(configListList, configListMap)
			}

			describeCloudNativeAPIGatewayConfigResultMap["config_list"] = []interface{}{configListList}
		}

		if result.GroupSubnetId != nil {
			describeCloudNativeAPIGatewayConfigResultMap["group_subnet_id"] = result.GroupSubnetId
		}

		if result.GroupVpcId != nil {
			describeCloudNativeAPIGatewayConfigResultMap["group_vpc_id"] = result.GroupVpcId
		}

		if result.GroupId != nil {
			describeCloudNativeAPIGatewayConfigResultMap["group_id"] = result.GroupId
		}

		ids = append(ids, *result.GatewayId)
		_ = d.Set("result", describeCloudNativeAPIGatewayConfigResultMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), describeCloudNativeAPIGatewayConfigResultMap); e != nil {
			return e
		}
	}
	return nil
}

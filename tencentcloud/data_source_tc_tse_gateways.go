/*
Use this data source to query detailed information of tse gateways

Example Usage

```hcl
data "tencentcloud_tse_gateways" "gateways" {
  filters {
		name = "Region"
		values =

  }
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

func dataSourceTencentCloudTseGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGatewaysRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions, valid value:Type,Name,GatewayId,Tag,TradeType,InternetPaymode,Region.",
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
							Description: "Filter value.",
						},
					},
				},
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Gateways information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total count.",
						},
						"gateway_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Gateway list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway ID.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of gateway. May return values: - Creating - CreateFailed - Running - Modifying - UpdatingSpec - UpdateFailed - Deleting - DeleteFailed - Isolating.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway name.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway type.",
									},
									"gateway_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway version. Reference value:- 2.4.1- 2.5.1.",
									},
									"node_config": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Original node config.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"specification": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specification, 1c2g|2c4g|4c8g|8c16g.",
												},
												"number": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Node number, 2-50.",
												},
											},
										},
									},
									"vpc_config": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Vpc infomation.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Subnet ID. Assign an IP address to the engine in the VPC subnet.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Subnet ID. Assign an IP address to the engine in the VPC subnet.",
												},
											},
										},
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of gateway.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Create time.",
									},
									"tags": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Tags infomation of gatewayNote: This field may return null, indicating that a valid value is not available.",
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
									"enable_cls": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to enable CLS log.",
									},
									"trade_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Trade type.- 0: postpaid- 1: Prepaid.",
									},
									"feature_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Product version.- TRIAL- STANDARD (default value)- PROFESSIONAL.",
									},
									"internet_max_bandwidth_out": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Public network outbound traffic bandwidth.",
									},
									"auto_renew_flag": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Auto renew flag- 0 ,default status- 1 ,auto renew- 2 ,auto not renew.",
									},
									"cur_deadline": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire date, for prepaid type.Note: This field may return null, indicating that a valid value is not available.",
									},
									"isolate_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Isolation time, used when the gateway is isolated.",
									},
									"enable_internet": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to open the public network of client.Note: This field may return null, indicating that a valid value is not available.",
									},
									"engine_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Engine region of gateway.",
									},
									"ingress_class_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ingress class name.",
									},
									"internet_pay_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Trade type of internet.- BANDWIDTH- TRAFFIC.",
									},
									"gateway_minor_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Minor version of gateway.",
									},
									"instance_port": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The port information that the instance monitors.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"http_port": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Http port.",
												},
												"https_port": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Https port.",
												},
											},
										},
									},
									"load_balancer_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Load balance type of public internet.",
									},
									"public_ip_addresses": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Addresses of public internet.",
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

func dataSourceTencentCloudTseGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tse_gateways.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tse.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tse.Filter{}
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

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*tse.ListCloudNativeAPIGatewayResult

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTseGatewaysByFilter(ctx, paramMap)
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
		listCloudNativeAPIGatewayResultMap := map[string]interface{}{}

		if result.TotalCount != nil {
			listCloudNativeAPIGatewayResultMap["total_count"] = result.TotalCount
		}

		if result.GatewayList != nil {
			gatewayListList := []interface{}{}
			for _, gatewayList := range result.GatewayList {
				gatewayListMap := map[string]interface{}{}

				if gatewayList.GatewayId != nil {
					gatewayListMap["gateway_id"] = gatewayList.GatewayId
				}

				if gatewayList.Status != nil {
					gatewayListMap["status"] = gatewayList.Status
				}

				if gatewayList.Name != nil {
					gatewayListMap["name"] = gatewayList.Name
				}

				if gatewayList.Type != nil {
					gatewayListMap["type"] = gatewayList.Type
				}

				if gatewayList.GatewayVersion != nil {
					gatewayListMap["gateway_version"] = gatewayList.GatewayVersion
				}

				if gatewayList.NodeConfig != nil {
					nodeConfigMap := map[string]interface{}{}

					if gatewayList.NodeConfig.Specification != nil {
						nodeConfigMap["specification"] = gatewayList.NodeConfig.Specification
					}

					if gatewayList.NodeConfig.Number != nil {
						nodeConfigMap["number"] = gatewayList.NodeConfig.Number
					}

					gatewayListMap["node_config"] = []interface{}{nodeConfigMap}
				}

				if gatewayList.VpcConfig != nil {
					vpcConfigMap := map[string]interface{}{}

					if gatewayList.VpcConfig.VpcId != nil {
						vpcConfigMap["vpc_id"] = gatewayList.VpcConfig.VpcId
					}

					if gatewayList.VpcConfig.SubnetId != nil {
						vpcConfigMap["subnet_id"] = gatewayList.VpcConfig.SubnetId
					}

					gatewayListMap["vpc_config"] = []interface{}{vpcConfigMap}
				}

				if gatewayList.Description != nil {
					gatewayListMap["description"] = gatewayList.Description
				}

				if gatewayList.CreateTime != nil {
					gatewayListMap["create_time"] = gatewayList.CreateTime
				}

				if gatewayList.Tags != nil {
					tagsList := []interface{}{}
					for _, tags := range gatewayList.Tags {
						tagsMap := map[string]interface{}{}

						if tags.TagKey != nil {
							tagsMap["tag_key"] = tags.TagKey
						}

						if tags.TagValue != nil {
							tagsMap["tag_value"] = tags.TagValue
						}

						tagsList = append(tagsList, tagsMap)
					}

					gatewayListMap["tags"] = []interface{}{tagsList}
				}

				if gatewayList.EnableCls != nil {
					gatewayListMap["enable_cls"] = gatewayList.EnableCls
				}

				if gatewayList.TradeType != nil {
					gatewayListMap["trade_type"] = gatewayList.TradeType
				}

				if gatewayList.FeatureVersion != nil {
					gatewayListMap["feature_version"] = gatewayList.FeatureVersion
				}

				if gatewayList.InternetMaxBandwidthOut != nil {
					gatewayListMap["internet_max_bandwidth_out"] = gatewayList.InternetMaxBandwidthOut
				}

				if gatewayList.AutoRenewFlag != nil {
					gatewayListMap["auto_renew_flag"] = gatewayList.AutoRenewFlag
				}

				if gatewayList.CurDeadline != nil {
					gatewayListMap["cur_deadline"] = gatewayList.CurDeadline
				}

				if gatewayList.IsolateTime != nil {
					gatewayListMap["isolate_time"] = gatewayList.IsolateTime
				}

				if gatewayList.EnableInternet != nil {
					gatewayListMap["enable_internet"] = gatewayList.EnableInternet
				}

				if gatewayList.EngineRegion != nil {
					gatewayListMap["engine_region"] = gatewayList.EngineRegion
				}

				if gatewayList.IngressClassName != nil {
					gatewayListMap["ingress_class_name"] = gatewayList.IngressClassName
				}

				if gatewayList.InternetPayMode != nil {
					gatewayListMap["internet_pay_mode"] = gatewayList.InternetPayMode
				}

				if gatewayList.GatewayMinorVersion != nil {
					gatewayListMap["gateway_minor_version"] = gatewayList.GatewayMinorVersion
				}

				if gatewayList.InstancePort != nil {
					instancePortMap := map[string]interface{}{}

					if gatewayList.InstancePort.HttpPort != nil {
						instancePortMap["http_port"] = gatewayList.InstancePort.HttpPort
					}

					if gatewayList.InstancePort.HttpsPort != nil {
						instancePortMap["https_port"] = gatewayList.InstancePort.HttpsPort
					}

					gatewayListMap["instance_port"] = []interface{}{instancePortMap}
				}

				if gatewayList.LoadBalancerType != nil {
					gatewayListMap["load_balancer_type"] = gatewayList.LoadBalancerType
				}

				if gatewayList.PublicIpAddresses != nil {
					gatewayListMap["public_ip_addresses"] = gatewayList.PublicIpAddresses
				}

				gatewayListList = append(gatewayListList, gatewayListMap)
			}

			listCloudNativeAPIGatewayResultMap["gateway_list"] = []interface{}{gatewayListList}
		}

		ids = append(ids, *result.GatewayId)
		_ = d.Set("result", listCloudNativeAPIGatewayResultMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), listCloudNativeAPIGatewayResultMap); e != nil {
			return e
		}
	}
	return nil
}

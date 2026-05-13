package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoMultiPathGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoMultiPathGatewaysRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. The maximum value of Filters.Values is 20. If this parameter is not filled in, all gateway information under the current appid will be returned. Detailed filter conditions are as follows: gateway-type: filter by gateway type, supporting values cloud and private, representing filtering cloud gateways and private gateways respectively; keyword: filter by gateway name keyword.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter field values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"gateways": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Gateway details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway ID.",
						},
						"gateway_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway name.",
						},
						"gateway_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway type. Valid values: cloud (cloud gateway managed by Tencent Cloud), private (private gateway deployed by user).",
						},
						"gateway_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Gateway port, range 1-65535 (excluding 8888).",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway status. Valid values: creating, online, offline, disable.",
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway IP, in IPv4 format.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway region ID, can be obtained from DescribeMultiPathGatewayRegions API.",
						},
						"lines": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Line information. Returned when querying gateway details via DescribeMultiPathGateway, not returned when querying gateway list via DescribeMultiPathGateways.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"line_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Line ID. line-0 and line-1 are built-in line IDs. line-0: direct line, does not support add/edit/delete; line-1: EdgeOne L4 proxy line, supports modifying instances and rules, does not support delete; line-2 and above: EdgeOne L4 proxy line or custom line, supports modify/delete instances and rules.",
									},
									"line_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Line type. Valid values: direct (direct line, does not support edit/delete), proxy (EdgeOne L4 proxy line, supports editing instances and rules, does not support delete), custom (custom line, supports edit and delete).",
									},
									"line_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Line address, in host:port format.",
									},
									"proxy_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "L4 proxy instance ID, returned when LineType is proxy (EdgeOne L4 proxy).",
									},
									"rule_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Forwarding rule ID, returned when LineType is proxy (EdgeOne L4 proxy).",
									},
								},
							},
						},
						"need_confirm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the gateway origin IP list has changed and needs re-confirmation. true: origin IP list changed, needs confirmation; false: origin IP list unchanged, no confirmation needed.",
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

func dataSourceTencentCloudTeoMultiPathGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_multi_path_gateways.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*teov20220901.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := teov20220901.Filter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					value := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(value))
				}
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var respData []*teov20220901.MultiPathGateway
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoMultiPathGatewaysByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	gatewaysList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, gateway := range respData {
			gatewayMap := map[string]interface{}{}
			if gateway.GatewayId != nil {
				gatewayMap["gateway_id"] = gateway.GatewayId
			}

			if gateway.GatewayName != nil {
				gatewayMap["gateway_name"] = gateway.GatewayName
			}

			if gateway.GatewayType != nil {
				gatewayMap["gateway_type"] = gateway.GatewayType
			}

			if gateway.GatewayPort != nil {
				gatewayMap["gateway_port"] = gateway.GatewayPort
			}

			if gateway.Status != nil {
				gatewayMap["status"] = gateway.Status
			}

			if gateway.GatewayIP != nil {
				gatewayMap["gateway_ip"] = gateway.GatewayIP
			}

			if gateway.RegionId != nil {
				gatewayMap["region_id"] = gateway.RegionId
			}

			if gateway.Lines != nil {
				linesList := make([]map[string]interface{}, 0, len(gateway.Lines))
				for _, line := range gateway.Lines {
					lineMap := map[string]interface{}{}
					if line.LineId != nil {
						lineMap["line_id"] = line.LineId
					}

					if line.LineType != nil {
						lineMap["line_type"] = line.LineType
					}

					if line.LineAddress != nil {
						lineMap["line_address"] = line.LineAddress
					}

					if line.ProxyId != nil {
						lineMap["proxy_id"] = line.ProxyId
					}

					if line.RuleId != nil {
						lineMap["rule_id"] = line.RuleId
					}

					linesList = append(linesList, lineMap)
				}

				gatewayMap["lines"] = linesList
			}

			if gateway.NeedConfirm != nil {
				gatewayMap["need_confirm"] = gateway.NeedConfirm
			}

			gatewaysList = append(gatewaysList, gatewayMap)
		}

		_ = d.Set("gateways", gatewaysList)
	}

	d.SetId(zoneId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), gatewaysList); e != nil {
			return e
		}
	}

	return nil
}

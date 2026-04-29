package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoMultiPathGateway() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceTencentCloudTeoMultiPathGatewayRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions for querying multi-path gateways. The detailed filtering conditions are as follows: <li>gateway-type: Filter by gateway type, supported values are cloud and private.</li><li>keyword: Filter by gateway name keyword.</li>.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Filter value.",
						},
					},
				},
			},

			// computed
			"gateways": {
				Computed:    true,
				Type:        schema.TypeList,
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
							Description: "Gateway type. Valid values: <li>cloud: Cloud gateway, created and managed by Tencent Cloud.</li><li>private: Private gateway, deployed by the user.</li>.",
						},
						"gateway_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Gateway port, range 1-65535 (excluding 8888).",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway status. Valid values: <li>creating: Creating;</li><li>online: Online;</li><li>offline: Offline;</li><li>disable: Disabled.</li>.",
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway IP, in IPv4 format.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway region ID, which can be obtained from the DescribeMultiPathGatewayRegions API.",
						},
						"lines": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Line information, returned when querying gateway details via DescribeMultiPathGateway, but not returned when querying gateway list via DescribeMultiPathGateways.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"line_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Line ID. line-0 and line-1 are built-in line IDs. Valid values: <li>line-0: Direct line, does not support adding, editing, or deletion;</li><li>line-1: EdgeOne Layer-4 proxy line, supports modifying instances and rules, does not support deletion;</li><li>line-2 and above: EdgeOne Layer-4 proxy line or custom line, supports modifying and deleting instances and rules.</li>.",
									},
									"line_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Line type. Valid values: <li>direct: Direct line, does not support editing or deletion;</li><li>proxy: EdgeOne Layer-4 proxy line, supports editing instances and rules, does not support deletion;</li><li>custom: Custom line, supports editing and deletion.</li>.",
									},
									"line_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Line address, in host:port format.",
									},
									"proxy_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Layer-4 proxy instance ID, returned when LineType is proxy (EdgeOne Layer-4 proxy).",
									},
									"rule_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Forwarding rule ID, returned when LineType is proxy (EdgeOne Layer-4 proxy).",
									},
								},
							},
						},
						"need_confirm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether reconfirmation is needed when the gateway origin IP list changes. Valid values: <li>true: The origin IP list has changed and needs confirmation;</li><li>false: The origin IP list has not changed and no confirmation is needed.</li>.",
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

func DataSourceTencentCloudTeoMultiPathGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_multi_path_gateway.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*teo.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filter := teo.Filter{}
			filterMap := item.(map[string]interface{})
			if v, ok := filterMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var gateways []*teo.MultiPathGateway
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoMultiPathGatewaysByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		gateways = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(gateways))
	tmpList := make([]map[string]interface{}, 0, len(gateways))
	if gateways != nil {
		for _, gateway := range gateways {
			gatewayMap := map[string]interface{}{}
			if gateway.GatewayId != nil {
				gatewayMap["gateway_id"] = gateway.GatewayId
				ids = append(ids, *gateway.GatewayId)
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

			tmpList = append(tmpList, gatewayMap)
		}

		_ = d.Set("gateways", tmpList)
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

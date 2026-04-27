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
				Description: "Multi-path gateway list.",
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
							Description: "Gateway type. Possible values are: cloud: cloud gateway; private: private gateway.",
						},
						"gateway_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Gateway port, range 1-65535 (excluding 8888).",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway status. Possible values are: creating: creating; online: online; offline: offline; disable: disabled.",
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway IP address in IPv4 format.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway region ID.",
						},
						"need_confirm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the gateway origin IP list has changed and needs re-confirmation. Possible values are: true: changed, needs confirmation; false: not changed, no confirmation needed.",
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

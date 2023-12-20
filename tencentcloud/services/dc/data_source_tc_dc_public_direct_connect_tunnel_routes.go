package dc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcPublicDirectConnectTunnelRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcPublicDirectConnectTunnelRoutesRead,
		Schema: map[string]*schema.Schema{
			"direct_connect_tunnel_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "direct connect tunnel id.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "filter condition: route-type: route type, value: BGP/STATIC route-subnet: route cidr, value such as: 192.68.1.0/24.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fields to be filtered.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "filter value of the field.",
						},
					},
				},
			},

			"routes": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Internet tunnel route list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "direct connect tunnel route id.",
						},
						"destination_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network CIDR.",
						},
						"route_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route type: BGP/STATIC route.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ENABLE: routing is enabled, DISABLE: routing is disabled.",
						},
						"as_path": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "ASPath info.",
						},
						"next_hop": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route next hop ip.",
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

func dataSourceTencentCloudDcPublicDirectConnectTunnelRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dc_public_direct_connect_tunnel_routes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("direct_connect_tunnel_id"); ok {
		paramMap["DirectConnectTunnelId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*dc.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := dc.Filter{}
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

	service := DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var routes []*dc.DirectConnectTunnelRoute

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcPublicDirectConnectTunnelRoutesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		routes = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(routes))
	tmpList := make([]map[string]interface{}, 0, len(routes))

	if routes != nil {
		for _, directConnectTunnelRoute := range routes {
			directConnectTunnelRouteMap := map[string]interface{}{}

			if directConnectTunnelRoute.RouteId != nil {
				directConnectTunnelRouteMap["route_id"] = directConnectTunnelRoute.RouteId
			}

			if directConnectTunnelRoute.DestinationCidrBlock != nil {
				directConnectTunnelRouteMap["destination_cidr_block"] = directConnectTunnelRoute.DestinationCidrBlock
			}

			if directConnectTunnelRoute.RouteType != nil {
				directConnectTunnelRouteMap["route_type"] = directConnectTunnelRoute.RouteType
			}

			if directConnectTunnelRoute.Status != nil {
				directConnectTunnelRouteMap["status"] = directConnectTunnelRoute.Status
			}

			if directConnectTunnelRoute.ASPath != nil {
				directConnectTunnelRouteMap["as_path"] = directConnectTunnelRoute.ASPath
			}

			if directConnectTunnelRoute.NextHop != nil {
				directConnectTunnelRouteMap["next_hop"] = directConnectTunnelRoute.NextHop
			}

			ids = append(ids, *directConnectTunnelRoute.RouteId)
			tmpList = append(tmpList, directConnectTunnelRouteMap)
		}

		_ = d.Set("routes", tmpList)
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

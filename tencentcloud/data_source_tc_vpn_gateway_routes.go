/*
Use this data source to query detailed information of VPN gateways.

Example Usage

```hcl
data "tencentcloud_vpn_gateways" "foo" {
  vpn_gateway_id              = "main"
  destination_cidr_block                = "vpngw-8ccsnclt"
  instance_type = "1.1.1.1"
  instance_id              = "ap-guangzhou-3"
  tags = {
    test = "tf"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpnGatewayRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpnGatewayRoutesRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPN gateway ID.",
			},
			"destination_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Destination IDC IP range.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Next hop type (type of the associated instance). Valid values: VPNCONN (VPN tunnel) and CCN (CCN instance).",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance ID of the next hop.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"vpn_gateway_route_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the vpn gateway routes.",
				Elem: &schema.Resource{
					Schema: VpnGatewayRoutePara(),
				},
			},
		},
	}
}

func dataSourceTencentCloudVpnGatewayRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpn_gateway_routes.read")()

	var (
		logId        = getLogId(contextNil)
		ctx          = context.WithValue(context.TODO(), logIdKey, logId)
		vpnGatewayId = d.Get("vpn_gateway_id").(string)
		vpcService   = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	params := make(map[string]string)
	if v, ok := d.GetOk("destination_cidr"); ok {
		params["DestinationCidr"] = v.(string)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		params["InstanceId"] = v.(string)
	}
	if v, ok := d.GetOk("instance_type"); ok {
		params["InstanceType"] = v.(string)
	}

	filters := make([]*vpc.Filter, 0, len(params))
	for k, v := range params {
		filter := &vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		filters = append(filters, filter)
	}
	err, result := vpcService.DescribeVpnGatewayRoutes(ctx, vpnGatewayId, filters)
	if err != nil {
		log.Printf("[CRITAL]%s read VPN gateway routes failed, reason:%s\n ", logId, err.Error())
		return err
	}
	ids := make([]string, 0, len(result))
	routeList := make([]map[string]interface{}, 0, len(result))
	for _, route := range result {
		routeList = append(routeList, ConverterVpnGatewayRouteToMap(vpnGatewayId, route))
		ids = append(ids, *route.RouteId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("vpn_gateway_route_list", routeList); e != nil {
		log.Printf("[CRITAL]%s provider set vpn gateway route list fail, reason:%s\n ", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), routeList); e != nil {
			return e
		}
	}

	return nil
}

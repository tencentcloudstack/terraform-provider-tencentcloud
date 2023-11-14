/*
Provides a resource to create a vpc vpn_gateway_ccn_routes

Example Usage

```hcl
resource "tencentcloud_vpc_vpn_gateway_ccn_routes" "vpn_gateway_ccn_routes" {
  vpn_gateway_id = "vpngw-c6orbuv7"
  routes {
		route_id = "vpnr-7t3tknmg"
		status = "ENABLE"
		destination_cidr_block = "10.2.2.0/24"

  }
}
```

Import

vpc vpn_gateway_ccn_routes can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_vpn_gateway_ccn_routes.vpn_gateway_ccn_routes vpn_gateway_ccn_routes_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudVpcVpnGatewayCcnRoutes() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcVpnGatewayCcnRoutesCreate,
		Read:   resourceTencentCloudVpcVpnGatewayCcnRoutesRead,
		Delete: resourceTencentCloudVpcVpnGatewayCcnRoutesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VPN GATEWAY INSTANCE ID.",
			},

			"routes": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "CCN Route(IDC network segment) List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Route Id.",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: " Whether routing information is enabled ENABLE：Enable Route DISABLE：Disable Route.",
						},
						"destination_cidr_block": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Routing CIDR.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpcVpnGatewayCcnRoutesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_gateway_ccn_routes.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = vpc.NewModifyVpnGatewayCcnRoutesRequest()
		response     = vpc.NewModifyVpnGatewayCcnRoutesResponse()
		vpnGatewayId string
		routeId      string
	)
	if v, ok := d.GetOk("vpn_gateway_id"); ok {
		vpnGatewayId = v.(string)
		request.VpnGatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("routes"); ok {
		for _, item := range v.([]interface{}) {
			vpngwCcnRoutes := vpc.VpngwCcnRoutes{}
			if v, ok := dMap["route_id"]; ok {
				vpngwCcnRoutes.RouteId = helper.String(v.(string))
			}
			if v, ok := dMap["status"]; ok {
				vpngwCcnRoutes.Status = helper.String(v.(string))
			}
			if v, ok := dMap["destination_cidr_block"]; ok {
				vpngwCcnRoutes.DestinationCidrBlock = helper.String(v.(string))
			}
			request.Routes = append(request.Routes, &vpngwCcnRoutes)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyVpnGatewayCcnRoutes(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc vpnGatewayCcnRoutes failed, reason:%+v", logId, err)
		return err
	}

	vpnGatewayId = *response.Response.VpnGatewayId
	d.SetId(strings.Join([]string{vpnGatewayId, routeId}, FILED_SP))

	return resourceTencentCloudVpcVpnGatewayCcnRoutesRead(d, meta)
}

func resourceTencentCloudVpcVpnGatewayCcnRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_gateway_ccn_routes.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcVpnGatewayCcnRoutesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_gateway_ccn_routes.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

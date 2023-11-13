/*
Provides a resource to create a vpc vpn_connection_reset

Example Usage

```hcl
resource "tencentcloud_vpc_vpn_connection_reset" "vpn_connection_reset" {
  vpn_gateway_id = "vpngw-c6orbuv7"
  vpn_connection_id = "vpnx-osftvdea"
}
```

Import

vpc vpn_connection_reset can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_vpn_connection_reset.vpn_connection_reset vpn_connection_reset_id
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

func resourceTencentCloudVpcVpnConnectionReset() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcVpnConnectionResetCreate,
		Read:   resourceTencentCloudVpcVpnConnectionResetRead,
		Delete: resourceTencentCloudVpcVpnConnectionResetDelete,
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

			"vpn_connection_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VPN CONNECTION INSTANCE ID.",
			},
		},
	}
}

func resourceTencentCloudVpcVpnConnectionResetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_connection_reset.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = vpc.NewResetVpnConnectionRequest()
		response        = vpc.NewResetVpnConnectionResponse()
		vpnGatewayId    string
		vpnConnectionId string
	)
	if v, ok := d.GetOk("vpn_gateway_id"); ok {
		vpnGatewayId = v.(string)
		request.VpnGatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpn_connection_id"); ok {
		vpnConnectionId = v.(string)
		request.VpnConnectionId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ResetVpnConnection(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc vpnConnectionReset failed, reason:%+v", logId, err)
		return err
	}

	vpnGatewayId = *response.Response.VpnGatewayId
	d.SetId(strings.Join([]string{vpnGatewayId, vpnConnectionId}, FILED_SP))

	return resourceTencentCloudVpcVpnConnectionResetRead(d, meta)
}

func resourceTencentCloudVpcVpnConnectionResetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_connection_reset.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcVpnConnectionResetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_connection_reset.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

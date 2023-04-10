/*
Provides a resource to create a vpc vpn_connection_reset

Example Usage

```hcl
resource "tencentcloud_vpn_connection_reset" "vpn_connection_reset" {
  vpn_gateway_id    = "vpngw-gt8bianl"
  vpn_connection_id = "vpnx-kme2tx8m"
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpnConnectionReset() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnConnectionResetCreate,
		Read:   resourceTencentCloudVpnConnectionResetRead,
		Delete: resourceTencentCloudVpnConnectionResetDelete,
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

func resourceTencentCloudVpnConnectionResetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpn_connection_reset.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = vpc.NewResetVpnConnectionRequest()
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
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc vpnConnectionReset failed, reason:%+v", logId, err)
		return nil
	}

	d.SetId(vpnGatewayId + FILED_SP + vpnConnectionId)

	return resourceTencentCloudVpnConnectionResetRead(d, meta)
}

func resourceTencentCloudVpnConnectionResetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_connection_reset.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpnConnectionResetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_connection_reset.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

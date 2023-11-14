/*
Provides a resource to create a vpc vpn_connection_default_health_check_ip_generate

Example Usage

```hcl
resource "tencentcloud_vpc_vpn_connection_default_health_check_ip_generate" "vpn_connection_default_health_check_ip_generate" {
  vpn_gateway_id = "vpngw-c6orbuv7"
}
```

Import

vpc vpn_connection_default_health_check_ip_generate can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_vpn_connection_default_health_check_ip_generate.vpn_connection_default_health_check_ip_generate vpn_connection_default_health_check_ip_generate_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudVpcVpnConnectionDefaultHealthCheckIpGenerate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcVpnConnectionDefaultHealthCheckIpGenerateCreate,
		Read:   resourceTencentCloudVpcVpnConnectionDefaultHealthCheckIpGenerateRead,
		Delete: resourceTencentCloudVpcVpnConnectionDefaultHealthCheckIpGenerateDelete,
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
		},
	}
}

func resourceTencentCloudVpcVpnConnectionDefaultHealthCheckIpGenerateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_connection_default_health_check_ip_generate.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = vpc.NewGenerateVpnConnectionDefaultHealthCheckIpRequest()
		response     = vpc.NewGenerateVpnConnectionDefaultHealthCheckIpResponse()
		vpnGatewayId string
	)
	if v, ok := d.GetOk("vpn_gateway_id"); ok {
		vpnGatewayId = v.(string)
		request.VpnGatewayId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().GenerateVpnConnectionDefaultHealthCheckIp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc vpnConnectionDefaultHealthCheckIpGenerate failed, reason:%+v", logId, err)
		return err
	}

	vpnGatewayId = *response.Response.VpnGatewayId
	d.SetId(vpnGatewayId)

	return resourceTencentCloudVpcVpnConnectionDefaultHealthCheckIpGenerateRead(d, meta)
}

func resourceTencentCloudVpcVpnConnectionDefaultHealthCheckIpGenerateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_connection_default_health_check_ip_generate.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcVpnConnectionDefaultHealthCheckIpGenerateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_connection_default_health_check_ip_generate.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

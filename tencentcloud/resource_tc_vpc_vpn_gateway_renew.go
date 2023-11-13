/*
Provides a resource to create a vpc vpn_gateway_renew

Example Usage

```hcl
resource "tencentcloud_vpc_vpn_gateway_renew" "vpn_gateway_renew" {
  vpn_gateway_id = "vpngw-c6orbuv7"
  instance_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_AUTO_RENEW"

  }
}
```

Import

vpc vpn_gateway_renew can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_vpn_gateway_renew.vpn_gateway_renew vpn_gateway_renew_id
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

func resourceTencentCloudVpcVpnGatewayRenew() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcVpnGatewayRenewCreate,
		Read:   resourceTencentCloudVpcVpnGatewayRenewRead,
		Delete: resourceTencentCloudVpcVpnGatewayRenewDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VPN Gateway Instance.",
			},

			"instance_charge_prepaid": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "PREPAID model.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The duration of purchasing an instance, unit: month. ranges:1,2,3,4,5,6,7,8,9,12,24,36.",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Automatic renew flag. values: NOTIFY_AND_AUTO_RENEW:Notify expiration and auto-renew.  NOTIFY_AND_MANUAL_RENEW:Notification expires without automatic renewal. Default value:NOTIFY_AND_MANUAL_RENEW.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpcVpnGatewayRenewCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_gateway_renew.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = vpc.NewRenewVpnGatewayRequest()
		response     = vpc.NewRenewVpnGatewayResponse()
		vpnGatewayId string
	)
	if v, ok := d.GetOk("vpn_gateway_id"); ok {
		vpnGatewayId = v.(string)
		request.VpnGatewayId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "instance_charge_prepaid"); ok {
		instanceChargePrepaid := vpc.InstanceChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			instanceChargePrepaid.Period = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok {
			instanceChargePrepaid.RenewFlag = helper.String(v.(string))
		}
		request.InstanceChargePrepaid = &instanceChargePrepaid
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().RenewVpnGateway(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc vpnGatewayRenew failed, reason:%+v", logId, err)
		return err
	}

	vpnGatewayId = *response.Response.VpnGatewayId
	d.SetId(vpnGatewayId)

	return resourceTencentCloudVpcVpnGatewayRenewRead(d, meta)
}

func resourceTencentCloudVpcVpnGatewayRenewRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_gateway_renew.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcVpnGatewayRenewDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_vpn_gateway_renew.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

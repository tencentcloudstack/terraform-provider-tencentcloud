/*
Provides a resource to create a vpc refresh_nat_dc_route

Example Usage

```hcl
resource "tencentcloud_nat_refresh_nat_dc_route" "refresh_nat_dc_route" {
  nat_gateway_id = "nat-gnxkey2e"
  vpc_id         = "vpc-pyyv5k3v"
  dry_run = true
}
```

Import

vpc refresh_nat_dc_route can be imported using the id, e.g.

```
terraform import tencentcloud_nat_refresh_nat_dc_route.refresh_nat_dc_route vpc_id#nat_gateway_id
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudNatRefreshNatDcRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudNatRefreshNatDcRouteCreate,
		Read:   resourceTencentCloudNatRefreshNatDcRouteRead,
		Delete: resourceTencentCloudNatRefreshNatDcRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unique identifier of Vpc.",
			},

			"nat_gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unique identifier of Nat Gateway.",
			},

			"dry_run": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to pre-refresh, valid values: True:yes, False:no.",
			},
		},
	}
}

func resourceTencentCloudNatRefreshNatDcRouteCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_nat_refresh_nat_dc_route.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = vpc.NewRefreshDirectConnectGatewayRouteToNatGatewayRequest()
		vpcId        string
		natGatewayId string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("nat_gateway_id"); ok {
		natGatewayId = v.(string)
		request.NatGatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request.DryRun = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().RefreshDirectConnectGatewayRouteToNatGateway(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc refreshNatDcRoute failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(vpcId + FILED_SP + natGatewayId)

	return resourceTencentCloudNatRefreshNatDcRouteRead(d, meta)
}

func resourceTencentCloudNatRefreshNatDcRouteRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_nat_refresh_nat_dc_route.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudNatRefreshNatDcRouteDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_nat_refresh_nat_dc_route.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

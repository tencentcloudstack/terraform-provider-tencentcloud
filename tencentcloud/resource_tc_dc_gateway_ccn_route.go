/*
Provides a resource to creating direct connect gateway route entry.

Example Usage

```hcl
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_dc_gateway" "ccn_main" {
  name                = "ci-cdg-ccn-test"
  network_instance_id = "${tencentcloud_ccn.main.id}"
  network_type        = "CCN"
  gateway_type        = "NORMAL"
}

resource "tencentcloud_dc_gateway_ccn_route" "route1" {
  dcg_id     = "${tencentcloud_dc_gateway.ccn_main.id}"
  cidr_block = "10.1.1.0/32"
}

resource "tencentcloud_dc_gateway_ccn_route" "route2" {
  dcg_id     = "${tencentcloud_dc_gateway.ccn_main.id}"
  cidr_block = "192.1.1.0/32"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func resourceTencentCloudDcGatewayCcnRouteInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcGatewayCcnRouteCreate,
		Read:   resourceTencentCloudDcGatewayCcnRouteRead,
		Delete: resourceTencentCloudDcGatewayCcnRouteDelete,
		Schema: map[string]*schema.Schema{
			"dcg_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the DCG",
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
				Description:  "A network address segment of IDC.",
			},

			//compute
			"as_path": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "As_Path list of the BGP.",
			},
		},
	}
}

func resourceTencentCloudDcGatewayCcnRouteCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_gateway_ccn_route.create")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		dcgId     = d.Get("dcg_id").(string)
		cidrBlock = d.Get("cidr_block").(string)
	)

	//Modification of this parameter[as_path] is not yet supported
	routeId, err := service.CreateDirectConnectGatewayCcnRoute(ctx, dcgId, cidrBlock, nil)

	if err != nil {
		return err
	}

	d.SetId(dcgId + "#" + routeId)

	return resourceTencentCloudDcGatewayCcnRouteRead(d, meta)
}

func resourceTencentCloudDcGatewayCcnRouteRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_gateway_ccn_route.read")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	items := strings.Split(d.Id(), "#")
	if len(items) != 2 {
		return fmt.Errorf("id of resource.tencentcloud_dc_gateway_ccn_route is wrong")
	}

	dcgId, routeId := items[0], items[1]
	info, has, err := service.DescribeDirectConnectGatewayCcnRoute(ctx, dcgId, routeId)
	if err != nil {
		return err
	}

	if has == 0 {
		d.SetId("")
		return nil
	}

	d.Set("dcg_id", info.dcgId)
	d.Set("cidr_block", info.cidrBlock)
	d.Set("as_path", info.asPaths)

	return nil
}

func resourceTencentCloudDcGatewayCcnRouteDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_gateway_ccn_route.delete")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	items := strings.Split(d.Id(), "#")
	if len(items) != 2 {
		return fmt.Errorf("id of resource.tencentcloud_dc_gateway_ccn_route is wrong")
	}

	dcgId, routeId := items[0], items[1]
	_, has, err := service.DescribeDirectConnectGatewayCcnRoute(ctx, dcgId, routeId)
	if err != nil {
		return err
	}

	if has == 0 {
		return nil
	}

	return service.DeleteDirectConnectGatewayCcnRoute(ctx, dcgId, routeId)
}

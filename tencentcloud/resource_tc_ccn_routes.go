/*
Provides a resource to create a vpc ccn_routes

Example Usage

```hcl
resource "tencentcloud_ccn_routes" "ccn_routes" {
  ccn_id = "ccn-39lqkygf"
  route_id = "ccnr-3o0dfyuw"
  switch = "on"
}
```

Import

vpc ccn_routes can be imported using the id, e.g.

```
terraform import tencentcloud_ccn_routes.ccn_routes ccnId#routesId
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func resourceTencentCloudCcnRoutes() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnRoutesCreate,
		Read:   resourceTencentCloudCcnRoutesRead,
		Update: resourceTencentCloudCcnRoutesUpdate,
		Delete: resourceTencentCloudCcnRoutesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "CCN Instance ID.",
			},

			"route_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "CCN Route Id List.",
			},

			"switch": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`on`: Enable, `off`: Disable.",
			},
		},
	}
}

func resourceTencentCloudCcnRoutesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_routes.create")()
	defer inconsistentCheck(d, meta)()

	ccnId := d.Get("ccn_id").(string)
	routeId := d.Get("route_id").(string)

	d.SetId(ccnId + FILED_SP + routeId)

	return resourceTencentCloudCcnRoutesUpdate(d, meta)
}

func resourceTencentCloudCcnRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_routes.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	ccnId := idSplit[0]
	routeId := idSplit[1]

	ccnRoutes, err := service.DescribeVpcCcnRoutesById(ctx, ccnId, routeId)
	if err != nil {
		return err
	}

	if ccnRoutes == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcCcnRoutes` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("ccn_id", ccnId)
	_ = d.Set("route_id", routeId)

	if ccnRoutes.Enabled != nil {
		if *ccnRoutes.Enabled {
			_ = d.Set("switch", "on")
		} else {
			_ = d.Set("switch", "off")
		}
	}

	return nil
}

func resourceTencentCloudCcnRoutesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_routes.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	ccnId := idSplit[0]
	routeId := idSplit[1]

	certSwitch := d.Get("switch").(string)

	if certSwitch == "on" {

		var (
			request = vpc.NewEnableCcnRoutesRequest()
		)
		request.CcnId = &ccnId
		request.RouteIds = []*string{&routeId}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().EnableCcnRoutes(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update vpc ccnRoutes failed, reason:%+v", logId, err)
			return err
		}
	} else {

		var (
			request = vpc.NewDisableCcnRoutesRequest()
		)
		request.CcnId = &ccnId
		request.RouteIds = []*string{&routeId}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DisableCcnRoutes(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update vpc ccnRoutes failed, reason:%+v", logId, err)
			return err
		}

	}

	return resourceTencentCloudCcnRoutesRead(d, meta)
}

func resourceTencentCloudCcnRoutesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_routes.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

/*
Provides a resource to create a mariadb renew_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_renew_instance" "renew_instance" {
  instance_id = "tdsql-9vqvls95"
  period      = 1
}
```
*/
package tencentcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbRenewInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbRenewInstanceCreate,
		Read:   resourceTencentCloudMariadbRenewInstanceRead,
		Delete: resourceTencentCloudMariadbRenewInstanceDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"period": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Renewal duration, unit: month.",
			},
		},
	}
}

func resourceTencentCloudMariadbRenewInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_renew_instance.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = mariadb.NewRenewDBInstanceRequest()
		instanceId string
		dealName   string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().RenewDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		dealName = *result.Response.DealName
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb renewInstance failed, reason:%+v", logId, err)
		return err
	}

	// check order
	OrderRequest := mariadb.NewDescribeOrdersRequest()
	OrderRequest.DealNames = common.StringPtrs([]string{dealName})
	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		resp, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().DescribeOrders(OrderRequest)
		if e != nil {
			return resource.RetryableError(err)
		}

		if resp == nil || resp.Response == nil {
			e = fmt.Errorf("TencentCloud SDK returns nil response, %s", request.GetAction())
			return resource.RetryableError(e)
		}

		if *resp.Response.TotalCount == 0 {
			e = fmt.Errorf("TencentCloud SDK returns empty deal")
			return resource.RetryableError(e)
		} else if len(resp.Response.Deals) > 1 {
			e = fmt.Errorf("TencentCloud SDK returns more than one deal")
			return resource.RetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb renewInstance task failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbRenewInstanceRead(d, meta)
}

func resourceTencentCloudMariadbRenewInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_renew_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbRenewInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_renew_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

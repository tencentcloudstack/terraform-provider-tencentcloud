/*
Provides a resource to create a mariadb close_db_extranet_access

Example Usage

```hcl
resource "tencentcloud_mariadb_close_db_extranet_access" "close_db_extranet_access" {
  instance_id = "tdsql-9vqvls95"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbCloseDBExtranetAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbCloseDBExtranetAccessCreate,
		Read:   resourceTencentCloudMariadbCloseDBExtranetAccessRead,
		Delete: resourceTencentCloudMariadbCloseDBExtranetAccessDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of instance for which to enable public network access. The ID is in the format of `tdsql-ow728lmc` and can be obtained through the `DescribeDBInstances` API.",
			},
			"ipv6_flag": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Whether IPv6 is used. Default value: 0.",
			},
		},
	}
}

func resourceTencentCloudMariadbCloseDBExtranetAccessCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_close_db_extranet_access.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = mariadb.NewCloseDBExtranetAccessRequest()
		instanceId string
		flowId     int64
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("ipv6_flag"); ok {
		request.Ipv6Flag = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CloseDBExtranetAccess(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb closeDBExtranetAccess failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeFlowById(ctx, flowId)
		if e != nil {
			return retryError(e)
		}

		if *result.Status == MARIADB_TASK_SUCCESS {
			return nil
		} else if *result.Status == MARIADB_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("operate mariadb closeDBExtranetAccess status is running"))
		} else if *result.Status == MARIADB_TASK_FAIL {
			return resource.NonRetryableError(fmt.Errorf("operate mariadb closeDBExtranetAccess status is fail"))
		} else {
			e = fmt.Errorf("operate mariadb closeDBExtranetAccess status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb closeDBExtranetAccess task failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbCloseDBExtranetAccessRead(d, meta)
}

func resourceTencentCloudMariadbCloseDBExtranetAccessRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_close_db_extranet_access.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbCloseDBExtranetAccessDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_close_db_extranet_access.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

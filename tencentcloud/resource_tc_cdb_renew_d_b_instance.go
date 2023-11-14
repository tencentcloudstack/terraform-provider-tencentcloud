/*
Provides a resource to create a cdb renew_d_b_instance

Example Usage

```hcl
resource "tencentcloud_cdb_renew_d_b_instance" "renew_d_b_instance" {
  instance_id = ""
  time_span =
  modify_pay_type = ""
}
```

Import

cdb renew_d_b_instance can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_renew_d_b_instance.renew_d_b_instance renew_d_b_instance_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCdbRenewDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbRenewDBInstanceCreate,
		Read:   resourceTencentCloudCdbRenewDBInstanceRead,
		Delete: resourceTencentCloudCdbRenewDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The instance ID to be renewed, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page, you can use [Query Instance List](https://cloud.tencent.com/document/api/236/ 15872).",
			},

			"time_span": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Renewal duration, unit: month, optional values include [1,2,3,4,5,6,7,8,9,10,11,12,24,36].",
			},

			"modify_pay_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "If you need to renew the Pay-As-You-Go instance to a Subscription instance, the value of this input parameter needs to be specified as `PREPAID`.",
			},
		},
	}
}

func resourceTencentCloudCdbRenewDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_renew_d_b_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cdb.NewRenewDBInstanceRequest()
		response   = cdb.NewRenewDBInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("time_span"); v != nil {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("modify_pay_type"); ok {
		request.ModifyPayType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().RenewDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cdb renewDBInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudCdbRenewDBInstanceRead(d, meta)
}

func resourceTencentCloudCdbRenewDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_renew_d_b_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCdbRenewDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_renew_d_b_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

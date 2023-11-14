/*
Provides a resource to create a redis renew_instance_operation

Example Usage

```hcl
resource "tencentcloud_redis_renew_instance_operation" "renew_instance_operation" {
  instance_id = "crs-c1nl9rpv"
  period = 1
  modify_pay_mode = "prepaid"
}
```

Import

redis renew_instance_operation can be imported using the id, e.g.

```
terraform import tencentcloud_redis_renew_instance_operation.renew_instance_operation renew_instance_operation_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudRedisRenewInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisRenewInstanceOperationCreate,
		Read:   resourceTencentCloudRedisRenewInstanceOperationRead,
		Delete: resourceTencentCloudRedisRenewInstanceOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"period": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Purchase duration, in months.",
			},

			"modify_pay_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Identifies whether the billing model is modified:The current instance billing mode is pay-as-you-go, which is prepaid and renewed.The billing mode of the current instance is subscription and you can not set this parameter.",
			},
		},
	}
}

func resourceTencentCloudRedisRenewInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_renew_instance_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = redis.NewRenewInstanceRequest()
		response   = redis.NewRenewInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("period"); v != nil {
		request.Period = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("modify_pay_mode"); ok {
		request.ModifyPayMode = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().RenewInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis renewInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudRedisRenewInstanceOperationRead(d, meta)
}

func resourceTencentCloudRedisRenewInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_renew_instance_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisRenewInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_renew_instance_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

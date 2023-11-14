/*
Provides a resource to create a redis clear_instance_operation

Example Usage

```hcl
resource "tencentcloud_redis_clear_instance_operation" "clear_instance_operation" {
  instance_id = "crs-c1nl9rpv"
  password = &lt;nil&gt;
}
```

Import

redis clear_instance_operation can be imported using the id, e.g.

```
terraform import tencentcloud_redis_clear_instance_operation.clear_instance_operation clear_instance_operation_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudRedisClearInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisClearInstanceOperationCreate,
		Read:   resourceTencentCloudRedisClearInstanceOperationRead,
		Delete: resourceTencentCloudRedisClearInstanceOperationDelete,
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

			"password": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Redis instance password (password-free instances do not need to pass passwords, non-password-free instances must be transmitted).",
			},
		},
	}
}

func resourceTencentCloudRedisClearInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_clear_instance_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = redis.NewClearInstanceRequest()
		response   = redis.NewClearInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ClearInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis clearInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 60*readRetryTimeout, time.Second, service.RedisClearInstanceOperationStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudRedisClearInstanceOperationRead(d, meta)
}

func resourceTencentCloudRedisClearInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_clear_instance_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisClearInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_clear_instance_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

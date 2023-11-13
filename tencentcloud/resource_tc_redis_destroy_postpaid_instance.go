/*
Provides a resource to create a redis destroy_postpaid_instance

Example Usage

```hcl
resource "tencentcloud_redis_destroy_postpaid_instance" "destroy_postpaid_instance" {
  instance_id = "crs-c1nl9rpv"
}
```

Import

redis destroy_postpaid_instance can be imported using the id, e.g.

```
terraform import tencentcloud_redis_destroy_postpaid_instance.destroy_postpaid_instance destroy_postpaid_instance_id
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

func resourceTencentCloudRedisDestroyPostpaidInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisDestroyPostpaidInstanceCreate,
		Read:   resourceTencentCloudRedisDestroyPostpaidInstanceRead,
		Delete: resourceTencentCloudRedisDestroyPostpaidInstanceDelete,
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
		},
	}
}

func resourceTencentCloudRedisDestroyPostpaidInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_destroy_postpaid_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = redis.NewDestroyPostpaidInstanceRequest()
		response   = redis.NewDestroyPostpaidInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().DestroyPostpaidInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis destroyPostpaidInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 60*readRetryTimeout, time.Second, service.RedisDestroyPostpaidInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 60*readRetryTimeout, time.Second, service.RedisDestroyPostpaidInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudRedisDestroyPostpaidInstanceRead(d, meta)
}

func resourceTencentCloudRedisDestroyPostpaidInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_destroy_postpaid_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisDestroyPostpaidInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_destroy_postpaid_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

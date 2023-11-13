/*
Provides a resource to create a redis read_only

Example Usage

```hcl
resource "tencentcloud_redis_read_only" "read_only" {
  instance_id = "crs-c1nl9rpv"
  input_mode = &lt;nil&gt;
}
```

Import

redis read_only can be imported using the id, e.g.

```
terraform import tencentcloud_redis_read_only.read_only read_only_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"log"
	"time"
)

func resourceTencentCloudRedisReadOnly() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisReadOnlyCreate,
		Read:   resourceTencentCloudRedisReadOnlyRead,
		Update: resourceTencentCloudRedisReadOnlyUpdate,
		Delete: resourceTencentCloudRedisReadOnlyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"input_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance input mode:0: read-write1: read-only.",
			},
		},
	}
}

func resourceTencentCloudRedisReadOnlyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_read_only.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisReadOnlyUpdate(d, meta)
}

func resourceTencentCloudRedisReadOnlyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_read_only.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	readOnlyId := d.Id()

	readOnly, err := service.DescribeRedisReadOnlyById(ctx, instanceId)
	if err != nil {
		return err
	}

	if readOnly == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisReadOnly` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if readOnly.InstanceId != nil {
		_ = d.Set("instance_id", readOnly.InstanceId)
	}

	if readOnly.InputMode != nil {
		_ = d.Set("input_mode", readOnly.InputMode)
	}

	return nil
}

func resourceTencentCloudRedisReadOnlyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_read_only.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := redis.NewModifyInstanceReadOnlyRequest()

	readOnlyId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "input_mode"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ModifyInstanceReadOnly(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis readOnly failed, reason:%+v", logId, err)
		return err
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 30*readRetryTimeout, time.Second, service.RedisReadOnlyStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudRedisReadOnlyRead(d, meta)
}

func resourceTencentCloudRedisReadOnlyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_read_only.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

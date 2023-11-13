/*
Provides a resource to create a redis connection_config

Example Usage

```hcl
resource "tencentcloud_redis_connection_config" "connection_config" {
  instance_id = "crs-c1nl9rpv"
  client_limit = &lt;nil&gt;
  bandwidth = &lt;nil&gt;
}
```

Import

redis connection_config can be imported using the id, e.g.

```
terraform import tencentcloud_redis_connection_config.connection_config connection_config_id
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

func resourceTencentCloudRedisConnectionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisConnectionConfigCreate,
		Read:   resourceTencentCloudRedisConnectionConfigRead,
		Update: resourceTencentCloudRedisConnectionConfigUpdate,
		Delete: resourceTencentCloudRedisConnectionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"client_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The total number of connections per shard.If read-only replicas are not enabled, the lower limit is 10,000 and the upper limit is 40,000.When you enable read-only replicas, the minimum limit is 10,000 and the upper limit is 10,000 × (the number of read replicas +3).",
			},

			"bandwidth": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Additional bandwidth, greater than 0, in MB/s.",
			},
		},
	}
}

func resourceTencentCloudRedisConnectionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_connection_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisConnectionConfigUpdate(d, meta)
}

func resourceTencentCloudRedisConnectionConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_connection_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	connectionConfigId := d.Id()

	connectionConfig, err := service.DescribeRedisConnectionConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if connectionConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisConnectionConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if connectionConfig.InstanceId != nil {
		_ = d.Set("instance_id", connectionConfig.InstanceId)
	}

	if connectionConfig.ClientLimit != nil {
		_ = d.Set("client_limit", connectionConfig.ClientLimit)
	}

	if connectionConfig.Bandwidth != nil {
		_ = d.Set("bandwidth", connectionConfig.Bandwidth)
	}

	return nil
}

func resourceTencentCloudRedisConnectionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_connection_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := redis.NewModifyConnectionConfigRequest()

	connectionConfigId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "client_limit", "bandwidth"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ModifyConnectionConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis connectionConfig failed, reason:%+v", logId, err)
		return err
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 30*readRetryTimeout, time.Second, service.RedisConnectionConfigStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudRedisConnectionConfigRead(d, meta)
}

func resourceTencentCloudRedisConnectionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_connection_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

/*
Provides a resource to create a redis ssl

Example Usage

```hcl
resource "tencentcloud_redis_ssl" "ssl" {
  instance_id = "crs-c1nl9rpv"
}
```

Import

redis ssl can be imported using the id, e.g.

```
terraform import tencentcloud_redis_ssl.ssl ssl_id
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

func resourceTencentCloudRedisSsl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisSslCreate,
		Read:   resourceTencentCloudRedisSslRead,
		Update: resourceTencentCloudRedisSslUpdate,
		Delete: resourceTencentCloudRedisSslDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},
		},
	}
}

func resourceTencentCloudRedisSslCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_ssl.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisSslUpdate(d, meta)
}

func resourceTencentCloudRedisSslRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_ssl.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	sslId := d.Id()

	ssl, err := service.DescribeRedisSslById(ctx, instanceId)
	if err != nil {
		return err
	}

	if ssl == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisSsl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ssl.InstanceId != nil {
		_ = d.Set("instance_id", ssl.InstanceId)
	}

	return nil
}

func resourceTencentCloudRedisSslUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_ssl.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		closeSSLRequest  = redis.NewCloseSSLRequest()
		closeSSLResponse = redis.NewCloseSSLResponse()
	)

	sslId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().CloseSSL(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis ssl failed, reason:%+v", logId, err)
		return err
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 30*readRetryTimeout, time.Second, service.RedisSslStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 30*readRetryTimeout, time.Second, service.RedisSslStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudRedisSslRead(d, meta)
}

func resourceTencentCloudRedisSslDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_ssl.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

/*
Provides a resource to create a redis connection_config

Example Usage

```hcl
resource "tencentcloud_redis_connection_config" "connection_config" {
  instance_id = "crs-c1nl9rpv"
  client_limit = "20000"
  bandwidth = "20"
}

```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisConnectionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisConnectionConfigCreate,
		Read:   resourceTencentCloudRedisConnectionConfigRead,
		Update: resourceTencentCloudRedisConnectionConfigUpdate,
		Delete: resourceTencentCloudRedisConnectionConfigDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
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
				Description: "Additional bandwidth, greater than 0, in MB.",
			},
		},
	}
}

func resourceTencentCloudRedisConnectionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_connection_config.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
	)
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

	instanceId := d.Id()

	connectionConfig, err := service.DescribeRedisInstanceById(ctx, instanceId)
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

	// if connectionConfig.Bandwidth != nil {
	// 	_ = d.Set("bandwidth", connectionConfig.Bandwidth)
	// }

	return nil
}

func resourceTencentCloudRedisConnectionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_connection_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := redis.NewModifyConnectionConfigRequest()
	response := redis.NewModifyConnectionConfigResponse()

	instanceId := d.Id()
	request.InstanceId = &instanceId

	if v, ok := d.GetOkExists("client_limit"); ok {
		request.ClientLimit = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("bandwidth"); ok {
		request.Bandwidth = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ModifyConnectionConfig(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "FailedOperation.SystemError" {
					return resource.NonRetryableError(e)
				}
			}
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis param failed, reason:%+v", logId, err)
		return err
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	taskId := *response.Response.TaskId
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeTaskInfo(ctx, instanceId, taskId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("change account is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis change connection fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisConnectionConfigRead(d, meta)
}

func resourceTencentCloudRedisConnectionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_connection_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

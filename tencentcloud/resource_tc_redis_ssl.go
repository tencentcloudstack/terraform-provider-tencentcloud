/*
Provides a resource to create a redis ssl

Example Usage

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 7
  region = "ap-guangzhou"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_redis_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[2].zone
  name              = "tf_redis_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_security_group" "foo" {
  name = "tf-redis-sg"
}

resource "tencentcloud_security_group_lite_rule" "foo" {
  security_group_id = tencentcloud_security_group.foo.id

  ingress = [
    "ACCEPT#192.168.1.0/24#80#TCP",
    "DROP#8.8.8.8#80,90#UDP",
    "DROP#0.0.0.0/0#80-90#TCP",
  ]

  egress = [
    "ACCEPT#192.168.0.0/16#ALL#TCP",
    "ACCEPT#10.0.0.0/8#ALL#ICMP",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]
}

resource "tencentcloud_redis_instance" "foo" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[2].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[2].type_id
  password           = "test12345789"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[2].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[2].redis_replicas_nums[0]
  name               = "terrform_test"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  security_groups    = [tencentcloud_security_group.foo.id]
}

resource "tencentcloud_redis_ssl" "ssl" {
  instance_id = tencentcloud_redis_instance.foo.id
  ssl_config = "disabled"
}
```

Import

redis ssl can be imported using the instanceId, e.g.

```
terraform import tencentcloud_redis_ssl.ssl crs-c1nl9rpv
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
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

			"ssl_config": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{"enabled", "disabled"}),
				Description:  "The SSL configuration status of the instance: `enabled`,`disabled`.",
			},
		},
	}
}

func resourceTencentCloudRedisSslCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_ssl.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
	)
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

	instanceId := d.Id()

	ssl, err := service.DescribeRedisSslById(ctx, instanceId)
	if err != nil {
		return err
	}

	if ssl == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisSsl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if *ssl.SSLConfig {
		_ = d.Set("ssl_config", "enabled")
	} else {
		_ = d.Set("ssl_config", "disabled")
	}

	return nil
}

func resourceTencentCloudRedisSslUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_ssl.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		openSSLRequest  = redis.NewOpenSSLRequest()
		closeSSLRequest = redis.NewCloseSSLRequest()
		taskId          int64
	)

	instanceId := d.Id()
	if v, ok := d.GetOkExists("ssl_config"); ok {
		config := v.(string)
		if config == "enabled" {
			openSSLRequest.InstanceId = &instanceId
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().OpenSSL(openSSLRequest)
				if e != nil {
					if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
						if ee.Code == "FailedOperation.SystemError" {
							return resource.NonRetryableError(e)
						}
					}
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, openSSLRequest.GetAction(), openSSLRequest.ToJsonString(), result.ToJsonString())
				}
				taskId = *result.Response.TaskId
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update redis ssl failed, reason:%+v", logId, err)
				return err
			}
		} else {
			closeSSLRequest.InstanceId = &instanceId
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().CloseSSL(closeSSLRequest)
				if e != nil {
					if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
						if ee.Code == "FailedOperation.SystemError" {
							return resource.NonRetryableError(e)
						}
					}
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, closeSSLRequest.GetAction(), closeSSLRequest.ToJsonString(), result.ToJsonString())
				}
				taskId = *result.Response.TaskId
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update redis ssl failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	if taskId > 0 {
		err := resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
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
				return resource.RetryableError(fmt.Errorf("ssl config is processing"))
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s redis ssl config fail, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudRedisSslRead(d, meta)
}

func resourceTencentCloudRedisSslDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_ssl.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

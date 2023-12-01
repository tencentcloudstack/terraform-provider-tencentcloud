/*
Provides a resource to create a redis replica_readonly

Example Usage

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 7
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_redis_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[1].zone
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
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[0].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[0].type_id
  password           = "test12345789"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[0].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[0].redis_replicas_nums[0]
  name               = "terrform_test"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  security_groups    = [tencentcloud_security_group.foo.id]
}

resource "tencentcloud_redis_replica_readonly" "replica_readonly" {
  instance_id = tencentcloud_redis_instance.foo.id
  readonly_policy = ["master"]
  operate = "enable"
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
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
)

func resourceTencentCloudRedisReplicaReadonly() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisReplicaReadonlyCreate,
		Read:   resourceTencentCloudRedisReplicaReadonlyRead,
		Update: resourceTencentCloudRedisReplicaReadonlyUpdate,
		Delete: resourceTencentCloudRedisReplicaReadonlyDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"readonly_policy": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Routing policy: Enter `master` or `replication`, which indicates the master node or slave node.",
			},

			"operate": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{"enable", "disable"}),
				Description:  "The replica is read-only, `enable` - enable read-write splitting, `disable`- disable read-write splitting.",
			},
		},
	}
}

func resourceTencentCloudRedisReplicaReadonlyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replica_readonly.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisReplicaReadonlyUpdate(d, meta)
}

func resourceTencentCloudRedisReplicaReadonlyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replica_readonly.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	instance, err := service.DescribeRedisInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisReplicaReadonly` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instance.InstanceId != nil {
		_ = d.Set("instance_id", instance.InstanceId)
	}

	if instance.SlaveReadWeight != nil {
		if *instance.SlaveReadWeight == 100 {
			_ = d.Set("operate", "enable")
		} else {
			_ = d.Set("operate", "disable")
		}
	}

	return nil
}

func resourceTencentCloudRedisReplicaReadonlyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replica_readonly.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		disableRequest  = redis.NewDisableReplicaReadonlyRequest()
		disableResponse = redis.NewDisableReplicaReadonlyResponse()
		enableRequest   = redis.NewEnableReplicaReadonlyRequest()
		enableResponse  = redis.NewEnableReplicaReadonlyResponse()
		taskId          int64
	)

	instanceId := d.Id()
	if v, ok := d.GetOk("operate"); ok {
		operate := v.(string)
		if operate == "enable" {
			enableRequest.InstanceId = &instanceId
			err := resource.Retry(3*writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().EnableReplicaReadonly(enableRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableRequest.GetAction(), enableRequest.ToJsonString(), result.ToJsonString())
				}
				enableResponse = result
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update redis replicaReadonly failed, reason:%+v", logId, err)
				return err
			}

			taskId = *enableResponse.Response.TaskId
		}
		if operate == "disable" {
			disableRequest.InstanceId = &instanceId
			err := resource.Retry(3*writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().DisableReplicaReadonly(disableRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableRequest.GetAction(), disableRequest.ToJsonString(), result.ToJsonString())
				}
				disableResponse = result
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update redis replicaReadonly failed, reason:%+v", logId, err)
				return err
			}

			taskId = *disableResponse.Response.TaskId
		}
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
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
			return resource.RetryableError(fmt.Errorf("change inputMode is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis change inputMode fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisReplicaReadonlyRead(d, meta)
}

func resourceTencentCloudRedisReplicaReadonlyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replica_readonly.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

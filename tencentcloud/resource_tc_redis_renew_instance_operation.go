/*
Provides a resource to create a redis renew_instance_operation

Example Usage

Renew Subscription Instances

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 7
}

data "tencentcloud_vpc" "vpc" {
  name = "Default-VPC"
}

data "tencentcloud_vpc_subnets" "subnet" {
  vpc_id = data.tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[1].zone
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
  vpc_id             = data.tencentcloud_vpc.vpc.id
  subnet_id          = data.tencentcloud_vpc_subnets.subnet.instance_list[0].subnet_id
  security_groups    = [tencentcloud_security_group.foo.id]
  charge_type        = "PREPAID"
  prepaid_period     = 1
}

resource "tencentcloud_redis_renew_instance_operation" "foo" {
  instance_id     = tencentcloud_redis_instance.foo.id
  period          = 1
  modify_pay_mode = "prepaid"
}
```

*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = redis.NewRenewInstanceRequest()
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
			if sdkErr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == "FailedOperation.PayFailed" {
					return resource.NonRetryableError(e)
				}
			}
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis renewInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	_, _, _, err = service.CheckRedisOnlineOk(ctx, instanceId, 20*readRetryTimeout)
	if err != nil {
		log.Printf("[CRITAL]%s redis upgradeVersionOperation fail, reason:%s\n", logId, err.Error())
		return err
	}

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

/*
Provides a resource to create a redis replicate_attachment

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

resource "tencentcloud_redis_instance" "instance" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[2].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[2].type_id
  password           = "test12345789"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[2].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[2].redis_replicas_nums[0]
  name               = "terrform_test_instance"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  security_groups    = [tencentcloud_security_group.foo.id]
}


resource "tencentcloud_redis_replicate_attachment" "replicate_attachment" {
  group_id = "crs-rpl-orfiwmn5"
  master_instance_id = tencentcloud_redis_instance.foo.id
  instance_ids = [tencentcloud_redis_instance.instance.id]
}
```

Import

redis replicate_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_redis_replicate_attachment.replicate_attachment replicate_attachment_id
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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisReplicateAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisReplicateAttachmentCreate,
		Read:   resourceTencentCloudRedisReplicateAttachmentRead,
		Update: resourceTencentCloudRedisReplicateAttachmentUpdate,
		Delete: resourceTencentCloudRedisReplicateAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of group.",
			},

			"master_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of master instance.",
			},

			"instance_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "All instance ids of the replication group.",
			},
		},
	}
}

func resourceTencentCloudRedisReplicateAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replicate_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request          = redis.NewAddReplicationInstanceRequest()
		groupId          string
		masterInstanceId string
		instanceIds      []string
	)

	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("master_instance_id"); ok {
		masterInstanceId = v.(string)
		instanceIds = append(instanceIds, masterInstanceId)
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceId := instanceIdsSet[i].(string)
			instanceIds = append(instanceIds, instanceId)
		}
	}

	d.SetId(groupId)

	for index, instanceId := range instanceIds {
		var instanceRole string
		if instanceId == masterInstanceId {
			if index == 0 {
				instanceRole = "rw"
			} else {
				continue
			}
		} else {
			instanceRole = "r"
		}
		service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
		err := service.AddReplicationInstance(ctx, groupId, instanceId, instanceRole)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudRedisReplicateAttachmentRead(d, meta)
}

func resourceTencentCloudRedisReplicateAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replicate_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	groupId := d.Id()

	replicateGroup, err := service.DescribeRedisReplicateInstanceById(ctx, groupId)
	if err != nil {
		return err
	}

	if replicateGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisReplicateAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("group_id", groupId)

	if replicateGroup.Instances != nil {
		instanceIds := make([]*string, 0)
		for _, v := range replicateGroup.Instances {
			if *v.Role == "rw" {
				_ = d.Set("master_instance_id", v.InstanceId)
			}
			instanceIds = append(instanceIds, v.InstanceId)
		}
		_ = d.Set("instance_ids", instanceIds)
	}

	return nil
}

func resourceTencentCloudRedisReplicateAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replicate_attachment.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request  = redis.NewChangeMasterInstanceRequest()
		response = redis.NewChangeMasterInstanceResponse()
	)

	groupId := d.Id()
	request.GroupId = &groupId

	if d.HasChange("master_instance_id") {
		instanceId := ""
		if v, ok := d.GetOk("master_instance_id"); ok {
			instanceId = v.(string)
			request.InstanceId = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ChangeMasterInstance(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate redis RedisReplicateAttachment failed, reason:%+v", logId, err)
			return err
		}

		taskId := *response.Response.TaskId

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
					return resource.RetryableError(fmt.Errorf("update redis changeMaster is processing"))
				}
			})

			if err != nil {
				log.Printf("[CRITAL]%s update redis changeMaster fail, reason:%s\n", logId, err.Error())
				return err
			}
		}
	}

	if d.HasChange("instance_ids") {
		oldInterface, newInterface := d.GetChange("instance_ids")
		oldInstances := oldInterface.(*schema.Set)
		newInstances := newInterface.(*schema.Set)
		remove := helper.InterfacesStrings(oldInstances.Difference(newInstances).List())
		add := helper.InterfacesStrings(newInstances.Difference(oldInstances).List())

		service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
		if len(add) > 0 {
			for _, instanceId := range add {
				err := service.AddReplicationInstance(ctx, groupId, instanceId, "r")
				if err != nil {
					return err
				}
			}
		}

		if len(remove) > 0 {
			for _, instanceId := range remove {
				err := service.DeleteRedisReplicateAttachmentById(ctx, instanceId, groupId)
				if err != nil {
					return err
				}
			}
		}
	}

	return resourceTencentCloudRedisReplicateAttachmentRead(d, meta)
}

func resourceTencentCloudRedisReplicateAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replicate_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	groupId := d.Id()

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	replicateGroup, err := service.DescribeRedisReplicateInstanceById(ctx, groupId)
	if err != nil {
		return err
	}

	if replicateGroup.Instances != nil {
		masterInstanceId := ""
		for _, v := range replicateGroup.Instances {
			if *v.Role == "rw" {
				masterInstanceId = *v.InstanceId
				continue
			}
			if err := service.DeleteRedisReplicateAttachmentById(ctx, *v.InstanceId, groupId); err != nil {
				return err
			}
		}
		if masterInstanceId != "" {
			if err := service.DeleteRedisReplicateAttachmentById(ctx, masterInstanceId, groupId); err != nil {
				return err
			}
		}
	}
	return nil
}

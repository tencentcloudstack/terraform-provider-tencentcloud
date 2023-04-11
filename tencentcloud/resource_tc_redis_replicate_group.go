/*
Provides a resource to create a redis replicate

Example Usage

```hcl
resource "tencentcloud_redis_replicate_group" "" {
  instance_id = "crs-c1nl9rpv"
  group_name = "group_1"
  remark = ""
}
```

Import

redis replicate group can be imported using the id, e.g.

```
terraform import tencentcloud_redis_replicate_group.replicate_group replicate_group_id
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

func resourceTencentCloudRedisReplicateGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisReplicateGroupCreate,
		Read:   resourceTencentCloudRedisReplicateGroupRead,
		Delete: resourceTencentCloudRedisReplicateGroupDelete,
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

			"group_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The group name.",
			},

			"remark": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Remark.",
			},
		},
	}
}

func resourceTencentCloudRedisReplicateGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replicate_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = redis.NewCreateReplicationGroupRequest()
		response   = redis.NewCreateReplicationGroupResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().CreateReplicationGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create redis replicate failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	taskId := *response.Response.TaskId
	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
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
			return resource.RetryableError(fmt.Errorf("create replicate is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis create replicate fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisReplicateGroupRead(d, meta)
}

func resourceTencentCloudRedisReplicateGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replicate_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	group, err := service.DescribeRedisReplicateById(ctx, instanceId)
	if err != nil {
		return err
	}

	if group == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisReplicate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if group.GroupName != nil {
		_ = d.Set("group_name", group.GroupName)
	}

	if group.Remark != nil {
		_ = d.Set("remark", group.Remark)
	}

	return nil
}

func resourceTencentCloudRedisReplicateGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replicate_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()

	taskId, err := service.DeleteRedisReplicateById(ctx, instanceId)
	if err != nil {
		return err
	}

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
			return resource.RetryableError(fmt.Errorf("delete replicate is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis delete replicate fail, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}

/*
Provides a resource to create a redis replicate_instance

Example Usage

```hcl
resource "tencentcloud_redis_replicate_instance" "replicate_instance" {
  instance_id = "crs-c1nl9rpv"
  group_id = "crs-rpl-c1nl9rpv"
  instance_role = "rw"
}
```

Import

redis replicate_instance can be imported using the id, e.g.

```
terraform import tencentcloud_redis_replicate_instance.replicate_instance replicate_instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisReplicateInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisReplicateInstanceCreate,
		Read:   resourceTencentCloudRedisReplicateInstanceRead,
		Update: resourceTencentCloudRedisReplicateInstanceUpdate,
		Delete: resourceTencentCloudRedisReplicateInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of group.",
			},

			"instance_role": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Assign roles to instances added to the replication group.:rw: read-write.r: read-only.",
			},
		},
	}
}

func resourceTencentCloudRedisReplicateInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replicate_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = redis.NewAddReplicationInstanceRequest()
		response   = redis.NewAddReplicationInstanceResponse()
		instanceId string
		groupId    string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_role"); ok {
		request.InstanceRole = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().AddReplicationInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create redis replicateInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + FILED_SP + groupId)

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
			return resource.RetryableError(fmt.Errorf("change param is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis change param fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisReplicateInstanceRead(d, meta)
}

func resourceTencentCloudRedisReplicateInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replicate_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	groupId := idSplit[1]

	replicateInstance, err := service.DescribeRedisReplicateInstanceById(ctx, instanceId, groupId)
	if err != nil {
		return err
	}

	if replicateInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisReplicateInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("group_id", groupId)

	if replicateInstance.Role != nil {
		_ = d.Set("instance_role", replicateInstance.Role)
	}

	return nil
}

func resourceTencentCloudRedisReplicateInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replicate_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		changeInstanceRoleRequest  = redis.NewChangeInstanceRoleRequest()
		changeInstanceRoleResponse = redis.NewChangeInstanceRoleResponse()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	groupId := idSplit[1]

	changeInstanceRoleRequest.InstanceId = &instanceId
	changeInstanceRoleRequest.GroupId = &groupId

	immutableArgs := []string{"instance_id", "group_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_role") {
		if v, ok := d.GetOk("instance_role"); ok {
			changeInstanceRoleRequest.InstanceRole = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ChangeInstanceRole(changeInstanceRoleRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, changeInstanceRoleRequest.GetAction(), changeInstanceRoleRequest.ToJsonString(), result.ToJsonString())
		}
		changeInstanceRoleResponse = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis replicateInstance failed, reason:%+v", logId, err)
		return err
	}

	taskId := *changeInstanceRoleResponse.Response.TaskId
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
			return resource.RetryableError(fmt.Errorf("change param is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis change param fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisReplicateInstanceRead(d, meta)
}

func resourceTencentCloudRedisReplicateInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replicate_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	groupId := idSplit[1]

	taskId, err := service.DeleteRedisReplicateInstanceById(ctx, instanceId, groupId)
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
			return resource.RetryableError(fmt.Errorf("change param is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis change param fail, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}

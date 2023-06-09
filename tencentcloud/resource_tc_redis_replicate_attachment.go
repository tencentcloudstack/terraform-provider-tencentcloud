/*
Provides a resource to create a redis replicate_attachment

Example Usage

```hcl
resource "tencentcloud_redis_replicate_attachment" "replicate_attachment" {
  instance_id = "crs-c1nl9rpv"
  group_id = "crs-rpl-c1nl9rpv"
  instance_role = "rw"
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
	"strings"

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
		Delete: resourceTencentCloudRedisReplicateAttachmentDelete,
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

			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of group.",
			},

			"instance_role": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Assign roles to instances added to the replication group.:rw: read-write.r: read-only.",
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
		log.Printf("[CRITAL]%s create redis replicateAttachment failed, reason:%+v", logId, err)
		return err
	}

	taskId := *response.Response.TaskId
	d.SetId(instanceId + FILED_SP + groupId)

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
				return resource.RetryableError(fmt.Errorf("Add replication is processing"))
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s redis add replication fail, reason:%s\n", logId, err.Error())
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

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	groupId := idSplit[1]

	replicateAttachment, err := service.DescribeRedisReplicateInstanceById(ctx, instanceId, groupId)
	if err != nil {
		return err
	}

	if replicateAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisReplicateAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("group_id", groupId)

	if replicateAttachment.Role != nil {
		_ = d.Set("instance_role", replicateAttachment.Role)
	}

	return nil
}

func resourceTencentCloudRedisReplicateAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replicate_attachment.delete")()
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

	taskId, err := service.DeleteRedisReplicateAttachmentById(ctx, instanceId, groupId)
	if err != nil {
		return err
	}

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
				return resource.RetryableError(fmt.Errorf("remove replication is processing"))
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s redis remove replication fail, reason:%s\n", logId, err.Error())
			return err
		}
	}
	return nil
}

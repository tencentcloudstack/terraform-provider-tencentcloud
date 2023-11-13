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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
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

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, groupId}, FILED_SP))

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 30*readRetryTimeout, time.Second, service.RedisReplicateAttachmentStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
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

	replicateAttachment, err := service.DescribeRedisReplicateAttachmentById(ctx, instanceId, groupId)
	if err != nil {
		return err
	}

	if replicateAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisReplicateAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if replicateAttachment.InstanceId != nil {
		_ = d.Set("instance_id", replicateAttachment.InstanceId)
	}

	if replicateAttachment.GroupId != nil {
		_ = d.Set("group_id", replicateAttachment.GroupId)
	}

	if replicateAttachment.InstanceRole != nil {
		_ = d.Set("instance_role", replicateAttachment.InstanceRole)
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

	if err := service.DeleteRedisReplicateAttachmentById(ctx, instanceId, groupId); err != nil {
		return err
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 30*readRetryTimeout, time.Second, service.RedisReplicateAttachmentStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

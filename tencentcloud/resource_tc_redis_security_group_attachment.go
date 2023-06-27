/*
Provides a resource to create a redis security_group_attachment

Example Usage

```hcl
resource "tencentcloud_redis_security_group_attachment" "security_group_attachment" {
  instance_id       = "crs-jf4ico4v"
  security_group_id = "sg-edmur627"
}
```

Import

redis security_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_redis_security_group_attachment.security_group_attachment instance_id#security_group_id
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
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

const PRODUCT string = "redis"

func resourceTencentCloudRedisSecurityGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisSecurityGroupAttachmentCreate,
		Read:   resourceTencentCloudRedisSecurityGroupAttachmentRead,
		Delete: resourceTencentCloudRedisSecurityGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"security_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Security group ID.",
			},
		},
	}
}

func resourceTencentCloudRedisSecurityGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_security_group_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = redis.NewAssociateSecurityGroupsRequest()
		securityGroupId string
		instanceId      string
	)
	request.Product = helper.String(PRODUCT)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = append(request.InstanceIds, &instanceId)
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		securityGroupId = v.(string)
		request.SecurityGroupId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().AssociateSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create redis securityGroupAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + FILED_SP + securityGroupId)

	return resourceTencentCloudRedisSecurityGroupAttachmentRead(d, meta)
}

func resourceTencentCloudRedisSecurityGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_security_group_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	securityGroupAttachment, err := service.DescribeRedisSecurityGroupAttachmentById(ctx, PRODUCT, instanceId, securityGroupId)
	if err != nil {
		return err
	}

	if securityGroupAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisSecurityGroupAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if securityGroupAttachment.SecurityGroupId != nil {
		_ = d.Set("security_group_id", securityGroupAttachment.SecurityGroupId)
	}

	return nil
}

func resourceTencentCloudRedisSecurityGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_security_group_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	if err := service.DeleteRedisSecurityGroupAttachmentById(ctx, PRODUCT, instanceId, securityGroupId); err != nil {
		return err
	}

	return nil
}

/*
Provides a resource to create a redis security_group_attachment

Example Usage

```hcl
resource "tencentcloud_redis_security_group_attachment" "security_group_attachment" {
  product = "redis"
  instance_ids =
  security_group_id = "crs-c1nl9rpv"
}
```

Import

redis security_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_redis_security_group_attachment.security_group_attachment security_group_attachment_id
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
)

func resourceTencentCloudRedisSecurityGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisSecurityGroupAttachmentCreate,
		Read:   resourceTencentCloudRedisSecurityGroupAttachmentRead,
		Delete: resourceTencentCloudRedisSecurityGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The name of the database engine, the value of this interface: redis.",
			},

			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of instance IDs, an array of one or more instance IDs.",
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
		response        = redis.NewAssociateSecurityGroupsResponse()
		securityGroupId string
		instanceId      string
	)
	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
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
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create redis securityGroupAttachment failed, reason:%+v", logId, err)
		return err
	}

	securityGroupId = *response.Response.SecurityGroupId
	d.SetId(strings.Join([]string{securityGroupId, instanceId}, FILED_SP))

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
	securityGroupId := idSplit[0]
	instanceId := idSplit[1]

	securityGroupAttachment, err := service.DescribeRedisSecurityGroupAttachmentById(ctx, securityGroupId, instanceId)
	if err != nil {
		return err
	}

	if securityGroupAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisSecurityGroupAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if securityGroupAttachment.Product != nil {
		_ = d.Set("product", securityGroupAttachment.Product)
	}

	if securityGroupAttachment.InstanceIds != nil {
		_ = d.Set("instance_ids", securityGroupAttachment.InstanceIds)
	}

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
	securityGroupId := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteRedisSecurityGroupAttachmentById(ctx, securityGroupId, instanceId); err != nil {
		return err
	}

	return nil
}

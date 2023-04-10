/*
Provides a resource to create a redis security_group

Example Usage

```hcl
resource "tencentcloud_redis_security_group" "security_group" {
  instance_id = "crs-c1nl9rpv"
  security_group_id = "sg-cyules4s5"
}
```

Import

redis security_group can be imported using the instance_id#security_group_id, e.g.

```
terraform import tencentcloud_redis_security_group.security_group crs-c1nl9rpv#sg-cyules4s5
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
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisSecurityGroupCreate,
		Read:   resourceTencentCloudRedisSecurityGroupRead,
		Delete: resourceTencentCloudRedisSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "The ID of instance.",
			},

			"security_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Security group ID.",
			},
		},
	}
}

func resourceTencentCloudRedisSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_security_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = redis.NewAssociateSecurityGroupsRequest()
		instanceId      string
		securityGroupId string
	)
	request.Product = helper.String("redis")

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = []*string{helper.String(v.(string))}
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
		log.Printf("[CRITAL]%s create redis securityGroup failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + FILED_SP + securityGroupId)
	return resourceTencentCloudRedisSecurityGroupRead(d, meta)
}

func resourceTencentCloudRedisSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_security_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]

	securityGroupId, err := service.DescribeRedisSecurityGroupById(ctx, instanceId)
	if err != nil {
		return err
	}

	if securityGroupId == "" {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisSecurityGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("security_group_id", securityGroupId)

	return nil
}

func resourceTencentCloudRedisSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_security_group.delete")()
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

	if err := service.DeleteRedisSecurityGroupById(ctx, instanceId, securityGroupId); err != nil {
		return err
	}

	return nil
}

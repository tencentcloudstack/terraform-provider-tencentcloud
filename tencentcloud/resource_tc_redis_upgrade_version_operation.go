/*
Provides a resource to create a redis upgrade_version_operation

Example Usage

```hcl
resource "tencentcloud_redis_upgrade_version_operation" "upgrade_version_operation" {
  instance_id = "crs-c1nl9rpv"
  target_instance_type = "6"
  switch_option = 2
}
```

*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisUpgradeVersionOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisUpgradeVersionOperationCreate,
		Read:   resourceTencentCloudRedisUpgradeVersionOperationRead,
		Delete: resourceTencentCloudRedisUpgradeVersionOperationDelete,
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

			"target_instance_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target instance type, same as [CreateInstances](https://cloud.tencent.com/document/api/239/20026), that is, the target type of the instance to be changed.",
			},

			"switch_option": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Switch mode:1 - maintenance time window switching,2 - immediate switching.",
			},
		},
	}
}

func resourceTencentCloudRedisUpgradeVersionOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade_version_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = redis.NewUpgradeInstanceVersionRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_instance_type"); ok {
		request.TargetInstanceType = helper.String(v.(string))
	}

	if v, _ := d.GetOk("switch_option"); v != nil {
		request.SwitchOption = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().UpgradeInstanceVersion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis upgradeVersionOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	_, _, _, err = service.CheckRedisOnlineOk(ctx, instanceId, 20*readRetryTimeout)
	if err != nil {
		log.Printf("[CRITAL]%s redis upgradeVersionOperation fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisUpgradeVersionOperationRead(d, meta)
}

func resourceTencentCloudRedisUpgradeVersionOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade_version_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisUpgradeVersionOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade_version_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

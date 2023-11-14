/*
Provides a resource to create a redis upgrade

Example Usage

```hcl
resource "tencentcloud_redis_upgrade" "upgrade" {
  instance_id = "crs-c1nl9rpv"
  start_time = "17:00"
  end_time = "19:00"
}
```

Import

redis upgrade can be imported using the id, e.g.

```
terraform import tencentcloud_redis_upgrade.upgrade upgrade_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"log"
)

func resourceTencentCloudRedisUpgrade() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisUpgradeCreate,
		Read:   resourceTencentCloudRedisUpgradeRead,
		Update: resourceTencentCloudRedisUpgradeUpdate,
		Delete: resourceTencentCloudRedisUpgradeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Maintenance window start time, e.g. 17:00.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The end time of the maintenance window, e.g. 19:00.",
			},
		},
	}
}

func resourceTencentCloudRedisUpgradeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisUpgradeUpdate(d, meta)
}

func resourceTencentCloudRedisUpgradeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	upgradeId := d.Id()

	upgrade, err := service.DescribeRedisUpgradeById(ctx, instanceId)
	if err != nil {
		return err
	}

	if upgrade == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisUpgrade` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if upgrade.InstanceId != nil {
		_ = d.Set("instance_id", upgrade.InstanceId)
	}

	if upgrade.StartTime != nil {
		_ = d.Set("start_time", upgrade.StartTime)
	}

	if upgrade.EndTime != nil {
		_ = d.Set("end_time", upgrade.EndTime)
	}

	return nil
}

func resourceTencentCloudRedisUpgradeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := redis.NewModifyMaintenanceWindowRequest()

	upgradeId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "start_time", "end_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ModifyMaintenanceWindow(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis upgrade failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRedisUpgradeRead(d, meta)
}

func resourceTencentCloudRedisUpgradeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_upgrade.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

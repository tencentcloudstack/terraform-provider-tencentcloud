/*
Provides a resource to create a redis backup_operation

Example Usage

```hcl
resource "tencentcloud_redis_backup_operation" "backup_operation" {
  instance_id = "crs-c1nl9rpv"
  remark = &lt;nil&gt;
  storage_days = 7
}
```

Import

redis backup_operation can be imported using the id, e.g.

```
terraform import tencentcloud_redis_backup_operation.backup_operation backup_operation_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudRedisBackupOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisBackupOperationCreate,
		Read:   resourceTencentCloudRedisBackupOperationRead,
		Delete: resourceTencentCloudRedisBackupOperationDelete,
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

			"remark": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Notes information for the backup.",
			},

			"storage_days": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Number of days to store.0 specifies the default retention time.",
			},
		},
	}
}

func resourceTencentCloudRedisBackupOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = redis.NewManualBackupInstanceRequest()
		response   = redis.NewManualBackupInstanceResponse()
		instanceId string
		backupId   string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, _ := d.GetOk("storage_days"); v != nil {
		request.StorageDays = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ManualBackupInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis backupOperation failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, backupId}, FILED_SP))

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 30*readRetryTimeout, time.Second, service.RedisBackupOperationStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudRedisBackupOperationRead(d, meta)
}

func resourceTencentCloudRedisBackupOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisBackupOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

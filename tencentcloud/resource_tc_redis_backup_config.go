/*
Use this data source to query which instance types of Redis are available in a specific region.

Example Usage

```hcl
resource "tencentcloud_redis_backup_config" "redislab" {
  redis_id       = "crs-7yl0q0dd"
  backup_time    = "04:00-05:00"
  backup_period  = ["Monday"]
}
```

Import

Redis  backup config can be imported, e.g.

```hcl
$ terraform import tencentcloud_redis_backup_config.redisconfig redis-id
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudRedisBackupConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisBackupConfigCreate,
		Read:   resourceTencentCloudRedisBackupConfigRead,
		Update: resourceTencentCloudRedisBackupConfigUpdate,
		Delete: resourceTencentCloudRedisBackupConfigDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"redis_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of a Redis instance to which the policy will be applied.",
			},
			"backup_time": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"00:00-01:00", "01:00-02:00", "02:00-03:00",
					"03:00-04:00", "04:00-05:00", "05:00-06:00",
					"06:00-07:00", "07:00-08:00", "08:00-09:00",
					"09:00-10:00", "10:00-11:00", "11:00-12:00",
					"12:00-13:00", "13:00-14:00", "14:00-15:00",
					"15:00-16:00", "16:00-17:00", "17:00-18:00",
					"18:00-19:00", "19:00-20:00", "20:00-21:00",
					"21:00-22:00", "22:00-23:00", "23:00-00:00",
				}),
				Description: "Specifys what time the backup action should take place.",
			},
			"backup_period": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifys which day the backup action should take place. Supported values include: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday and Sunday.",
			},
		},
	}
}

func resourceTencentCloudRedisBackupConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup_config.create")()

	d.SetId(d.Get("redis_id").(string))

	return resourceTencentCloudRedisBackupConfigUpdate(d, meta)
}

func resourceTencentCloudRedisBackupConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup_config.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		backupPeriods, backupTime, e := service.DescribeAutoBackupConfig(ctx, d.Id())
		if e != nil {
			if sdkErr, ok := e.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == RedisInstanceNotFound {
					d.SetId("")
					return nil
				}
			}
			return retryError(e)
		}
		_ = d.Set("backup_time", backupTime)
		_ = d.Set("redis_id", d.Id())
		_ = d.Set("backup_period", backupPeriods)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudRedisBackupConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup_config.update")()

	weeksAllows := map[string]bool{
		"Monday": true, "Tuesday": true, "Wednesday": true,
		"Thursday": true, "Friday": true, "Saturday": true, "Sunday": true,
	}

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		isUpdate   = false
		redisId    = d.Get("redis_id").(string)
		backupTime = d.Get("backup_time").(string)
	)

	interfaceBackupPeriods := d.Get("backup_period").(*schema.Set).List()
	backupPeriods := make([]string, 0, len(interfaceBackupPeriods))

	for _, v := range interfaceBackupPeriods {
		if !weeksAllows[v.(string)] {
			return fmt.Errorf("redis backup config[backup_period] not supports %s", v.(string))
		}
		backupPeriods = append(backupPeriods, v.(string))
	}
	if len(backupPeriods) == 0 {
		return fmt.Errorf("redis backup config[backup_period] can not empty")
	}

	if d.HasChange("backup_time") || d.HasChange("backup_period") {
		isUpdate = true
	}

	if isUpdate {
		err := service.ModifyAutoBackupConfig(ctx, redisId, backupPeriods, backupTime)
		if err != nil {
			return err
		}
	}
	return resourceTencentCloudRedisBackupConfigRead(d, meta)
}

func resourceTencentCloudRedisBackupConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_backup_config.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		backupPeriods = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		backupTime    = "00:00-01:00"
	)

	err := service.ModifyAutoBackupConfig(ctx, d.Id(), backupPeriods, backupTime)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

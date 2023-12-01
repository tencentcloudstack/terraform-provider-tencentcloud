/*
Use this resource to create a backup config of redis.

Example Usage

Set configuration for automatic backups

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 7
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_redis_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[1].zone
  name              = "tf_redis_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_redis_instance" "foo" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[1].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[1].type_id
  password           = "test12345789"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[1].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[1].redis_replicas_nums[0]
  name               = "terrform_test"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_redis_backup_config" "foo" {
  redis_id      = tencentcloud_redis_instance.foo.id
  backup_time   = "04:00-05:00"
  backup_period = ["Monday"]
}
```

Import

Redis  backup config can be imported, e.g.

```
$ terraform import tencentcloud_redis_backup_config.foo redis-id
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Description: "ID of a redis instance to which the policy will be applied.",
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
				Description: "Specifys what time the backup action should take place. And the time interval should be one hour.",
			},
			"backup_period": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Deprecated:  "It has been deprecated from version 1.58.2. It makes no difference to online config at all",
				Description: "Specifys which day the backup action should take place. Valid values: `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday` and `Sunday`.",
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
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		isUpdate   = false
		redisId    = d.Get("redis_id").(string)
		backupTime = d.Get("backup_time").(string)
	)

	interfaceBackupPeriods := d.Get("backup_period").(*schema.Set).List()

	for _, v := range interfaceBackupPeriods {
		if !weeksAllows[v.(string)] {
			return fmt.Errorf("redis backup config[backup_period] not supports %s", v.(string))
		}
	}

	backupPeriods := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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

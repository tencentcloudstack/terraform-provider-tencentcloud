/*
Provides a resource to create a postgres backup_plan_config

Example Usage

```hcl
resource "tencentcloud_postgresql_backup_plan_config" "backup_plan_config" {
  db_instance_id = local.pgsql_id
  min_backup_start_time = "01:00:00"
  max_backup_start_time = "02:00:00"
  base_backup_retention_period = 7
  backup_period = ["monday","wednesday","friday"]
}
```

Import

postgres backup_plan_config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_backup_plan_config.backup_plan_config backup_plan_config_id
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlBackupPlanConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlBackupPlanConfigCreate,
		Read:   resourceTencentCloudPostgresqlBackupPlanConfigRead,
		Update: resourceTencentCloudPostgresqlBackupPlanConfigUpdate,
		Delete: resourceTencentCloudPostgresqlBackupPlanConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"min_backup_start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The earliest time to start a backup.",
			},

			"max_backup_start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The latest time to start a backup.",
			},

			"base_backup_retention_period": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Backup retention period in days. Value range:3-7.",
			},

			"backup_period": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Backup cycle, which means on which days each week the instance will be backed up. The parameter value should be the lowercase names of the days of the week.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlBackupPlanConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_backup_plan_config.create")()
	defer inconsistentCheck(d, meta)()

	var dBInstanceId string
	if v, ok := d.GetOk("db_instance_id"); ok {
		dBInstanceId = v.(string)
	}

	d.SetId(dBInstanceId)

	return resourceTencentCloudPostgresqlBackupPlanConfigUpdate(d, meta)
}

func resourceTencentCloudPostgresqlBackupPlanConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_backup_plan_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	dBInstanceId := d.Id()

	BackupPlanConfig, err := service.DescribePostgresqlBackupPlanConfigById(ctx, dBInstanceId)
	if err != nil {
		return err
	}

	if BackupPlanConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresBackupPlanConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("db_instance_id", dBInstanceId)

	if BackupPlanConfig.MinBackupStartTime != nil {
		_ = d.Set("min_backup_start_time", BackupPlanConfig.MinBackupStartTime)
	}

	if BackupPlanConfig.MaxBackupStartTime != nil {
		_ = d.Set("max_backup_start_time", BackupPlanConfig.MaxBackupStartTime)
	}

	if BackupPlanConfig.BaseBackupRetentionPeriod != nil {
		_ = d.Set("base_backup_retention_period", BackupPlanConfig.BaseBackupRetentionPeriod)
	}

	if BackupPlanConfig.BackupPeriod != nil {
		var newJson interface{}
		err := json.Unmarshal([]byte(*BackupPlanConfig.BackupPeriod), &newJson)
		if err != nil {
			return fmt.Errorf("convert BackupPeriod from string to interface{} failed, reason:+%v", err.Error())
		}
		_ = d.Set("backup_period", newJson)

	}

	return nil
}

func resourceTencentCloudPostgresqlBackupPlanConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_backup_plan_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = postgres.NewModifyBackupPlanRequest()
	)

	if d.HasChange("db_instance_id") {
		if v, ok := d.GetOk("db_instance_id"); ok {
			request.DBInstanceId = helper.String(v.(string))
		}
	}

	if d.HasChange("min_backup_start_time") {
		if v, ok := d.GetOk("min_backup_start_time"); ok {
			request.MinBackupStartTime = helper.String(v.(string))
		}
	}

	if d.HasChange("max_backup_start_time") {
		if v, ok := d.GetOk("max_backup_start_time"); ok {
			request.MaxBackupStartTime = helper.String(v.(string))
		}
	}

	if d.HasChange("base_backup_retention_period") {
		if v, ok := d.GetOkExists("base_backup_retention_period"); ok {
			request.BaseBackupRetentionPeriod = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("backup_period") {
		if v, ok := d.GetOk("backup_period"); ok {
			backupPeriodSet := v.(*schema.Set).List()
			for i := range backupPeriodSet {
				if backupPeriodSet[i] != nil {
					backupPeriod := backupPeriodSet[i].(string)
					request.BackupPeriod = append(request.BackupPeriod, &backupPeriod)
				}
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().ModifyBackupPlan(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update postgres BackupPlanConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresqlBackupPlanConfigRead(d, meta)
}

func resourceTencentCloudPostgresqlBackupPlanConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_backup_plan_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

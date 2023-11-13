/*
Provides a resource to create a postgres backup_plan

Example Usage

```hcl
resource "tencentcloud_postgres_backup_plan" "backup_plan" {
  d_b_instance_id = "postgres-xxxxx"
  min_backup_start_time = "01:00:00"
  max_backup_start_time = "02:00:00"
  base_backup_retention_period = 7
  backup_period =
}
```

Import

postgres backup_plan can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_backup_plan.backup_plan backup_plan_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"log"
)

func resourceTencentCloudPostgresBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresBackupPlanCreate,
		Read:   resourceTencentCloudPostgresBackupPlanRead,
		Update: resourceTencentCloudPostgresBackupPlanUpdate,
		Delete: resourceTencentCloudPostgresBackupPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
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

func resourceTencentCloudPostgresBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_backup_plan.create")()
	defer inconsistentCheck(d, meta)()

	var dBInstanceId string
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		dBInstanceId = v.(string)
	}

	d.SetId(dBInstanceId)

	return resourceTencentCloudPostgresBackupPlanUpdate(d, meta)
}

func resourceTencentCloudPostgresBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_backup_plan.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	backupPlanId := d.Id()

	BackupPlan, err := service.DescribePostgresBackupPlanById(ctx, dBInstanceId)
	if err != nil {
		return err
	}

	if BackupPlan == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresBackupPlan` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if BackupPlan.DBInstanceId != nil {
		_ = d.Set("d_b_instance_id", BackupPlan.DBInstanceId)
	}

	if BackupPlan.MinBackupStartTime != nil {
		_ = d.Set("min_backup_start_time", BackupPlan.MinBackupStartTime)
	}

	if BackupPlan.MaxBackupStartTime != nil {
		_ = d.Set("max_backup_start_time", BackupPlan.MaxBackupStartTime)
	}

	if BackupPlan.BaseBackupRetentionPeriod != nil {
		_ = d.Set("base_backup_retention_period", BackupPlan.BaseBackupRetentionPeriod)
	}

	if BackupPlan.BackupPeriod != nil {
		_ = d.Set("backup_period", BackupPlan.BackupPeriod)
	}

	return nil
}

func resourceTencentCloudPostgresBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_backup_plan.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := postgres.NewModifyBackupPlanRequest()

	backupPlanId := d.Id()

	request.DBInstanceId = &dBInstanceId

	immutableArgs := []string{"d_b_instance_id", "min_backup_start_time", "max_backup_start_time", "base_backup_retention_period", "backup_period"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().ModifyBackupPlan(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update postgres BackupPlan failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresBackupPlanRead(d, meta)
}

func resourceTencentCloudPostgresBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_backup_plan.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

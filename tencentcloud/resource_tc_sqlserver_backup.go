/*
Provides a resource to create a sqlserver backup

Example Usage

```hcl
resource "tencentcloud_sqlserver_backup" "backup" {
  strategy =
  db_names =
  instance_id = ""
  backup_name = ""
  backup_id = ""
  group_id = ""
}
```

Import

sqlserver backup can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_backup.backup backup_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverBackupCreate,
		Read:   resourceTencentCloudSqlserverBackupRead,
		Update: resourceTencentCloudSqlserverBackupUpdate,
		Delete: resourceTencentCloudSqlserverBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"strategy": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Backup strategy (0-instance backup 1-multi-database backup).",
			},

			"db_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of database names to be backed up (fill in only for multi-database backup).",
			},

			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, such as mssql-i1z41iwd.",
			},

			"backup_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The backup name, if not filled in, will automatically generate &amp;quot;Instance ID_Backup Start Time Stamp&amp;quot;.",
			},

			"backup_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The backup ID to be modified. It can be obtained through the `DescribeBackups` interface.",
			},

			"group_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The group ID of backup task. In the single-database backup file mode, it can be obtained through the `DescribeBackups` interface. BackupId and GroupId exist at the same time, modify according to BackupId.",
			},
		},
	}
}

func resourceTencentCloudSqlserverBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_backup.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = sqlserver.NewCreateBackupRequest()
		response   = sqlserver.NewCreateBackupResponse()
		instanceId string
	)
	if v, _ := d.GetOk("strategy"); v != nil {
		request.Strategy = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("db_names"); ok {
		dbNamesSet := v.(*schema.Set).List()
		for i := range dbNamesSet {
			dbName := dbNamesSet[i].(string)
			request.DBNames = append(request.DBNames, helper.String(dbName))
		}
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_name"); ok {
		request.BackupName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CreateBackup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver backup failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.FlowId
	d.SetId(helper.String(instanceId))

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"1"}, 2*readRetryTimeout, time.Second, service.SqlserverBackupStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudSqlserverBackupRead(d, meta)
}

func resourceTencentCloudSqlserverBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_backup.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	backupId := d.Id()

	backup, err := service.DescribeSqlserverBackupById(ctx, instanceId)
	if err != nil {
		return err
	}

	if backup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverBackup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if backup.Strategy != nil {
		_ = d.Set("strategy", backup.Strategy)
	}

	if backup.DbNames != nil {
		_ = d.Set("db_names", backup.DbNames)
	}

	if backup.InstanceId != nil {
		_ = d.Set("instance_id", backup.InstanceId)
	}

	if backup.BackupName != nil {
		_ = d.Set("backup_name", backup.BackupName)
	}

	if backup.BackupId != nil {
		_ = d.Set("backup_id", backup.BackupId)
	}

	if backup.GroupId != nil {
		_ = d.Set("group_id", backup.GroupId)
	}

	return nil
}

func resourceTencentCloudSqlserverBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_backup.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyBackupNameRequest()

	backupId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"strategy", "db_names", "instance_id", "backup_name", "backup_id", "group_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("backup_name") {
		if v, ok := d.GetOk("backup_name"); ok {
			request.BackupName = helper.String(v.(string))
		}
	}

	if d.HasChange("backup_id") {
		if v, ok := d.GetOk("backup_id"); ok {
			request.BackupId = helper.String(v.(string))
		}
	}

	if d.HasChange("group_id") {
		if v, ok := d.GetOk("group_id"); ok {
			request.GroupId = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyBackupName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver backup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverBackupRead(d, meta)
}

func resourceTencentCloudSqlserverBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_backup.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	backupId := d.Id()

	if err := service.DeleteSqlserverBackupById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}

/*
Provides a resource to create a sqlserver general_backups

Example Usage

```hcl
resource "tencentcloud_sqlserver_general_backups" "general_backups" {
  strategy = 0
  d_b_names =
  instance_id = "mssql-i1z41iwd"
  backup_name = "bk_name"
}
```

Import

sqlserver general_backups can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_backups.general_backups general_backups_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSqlserverGeneralBackups() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralBackupsCreate,
		Read:   resourceTencentCloudSqlserverGeneralBackupsRead,
		Update: resourceTencentCloudSqlserverGeneralBackupsUpdate,
		Delete: resourceTencentCloudSqlserverGeneralBackupsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"strategy": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Backup policy (0: instance backup, 1: multi-database backup).",
			},

			"d_b_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of names of databases to be backed up (required only for multi-database backup).",
			},

			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of mssql-i1z41iwd.",
			},

			"backup_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Backup name. If this parameter is left empty, a backup name in the format of [Instance ID]_[Backup start timestamp] will be automatically generated.",
			},
		},
	}
}

func resourceTencentCloudSqlserverGeneralBackupsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_backups.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = sqlserver.NewCreateBackupRequest()
		response   = sqlserver.NewCreateBackupResponse()
		instanceId string
	)
	if v, ok := d.GetOkExists("strategy"); ok {
		request.Strategy = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("d_b_names"); ok {
		dBNamesSet := v.(*schema.Set).List()
		for i := range dBNamesSet {
			dBNames := dBNamesSet[i].(string)
			request.DBNames = append(request.DBNames, &dBNames)
		}
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
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
		log.Printf("[CRITAL]%s create sqlserver generalBackups failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudSqlserverGeneralBackupsRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralBackupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_backups.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	generalBackupsId := d.Id()

	generalBackups, err := service.DescribeSqlserverGeneralBackupsById(ctx, instanceId)
	if err != nil {
		return err
	}

	if generalBackups == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverGeneralBackups` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if generalBackups.Strategy != nil {
		_ = d.Set("strategy", generalBackups.Strategy)
	}

	if generalBackups.DBNames != nil {
		_ = d.Set("d_b_names", generalBackups.DBNames)
	}

	if generalBackups.InstanceId != nil {
		_ = d.Set("instance_id", generalBackups.InstanceId)
	}

	if generalBackups.BackupName != nil {
		_ = d.Set("backup_name", generalBackups.BackupName)
	}

	return nil
}

func resourceTencentCloudSqlserverGeneralBackupsUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_backups.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyBackupNameRequest()

	generalBackupsId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"strategy", "d_b_names", "instance_id", "backup_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_id") {
		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}
	}

	if d.HasChange("backup_name") {
		if v, ok := d.GetOk("backup_name"); ok {
			request.BackupName = helper.String(v.(string))
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
		log.Printf("[CRITAL]%s update sqlserver generalBackups failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverGeneralBackupsRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralBackupsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_backups.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	generalBackupsId := d.Id()

	if err := service.DeleteSqlserverGeneralBackupsById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}

/*
Provides a resource to create a sqlserver restore_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_restore_instance" "restore_instance" {
  instance_id = "mssql-i1z41iwd"
  backup_id = 1981910
  target_instance_id = "mssql-au8ajamz"
  rename_restore {
		old_name = ""
		new_name = ""

  }
  type =
  d_b_list =
  group_id = ""
}
```

Import

sqlserver restore_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_restore_instance.restore_instance restore_instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"log"
)

func resourceTencentCloudSqlserverRestoreInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverRestoreInstanceCreate,
		Read:   resourceTencentCloudSqlserverRestoreInstanceRead,
		Update: resourceTencentCloudSqlserverRestoreInstanceUpdate,
		Delete: resourceTencentCloudSqlserverRestoreInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"backup_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Backup file ID, which can be obtained through the Id field in the returned value of the DescribeBackups API.",
			},

			"target_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of the target instance to which the backup is restored. The target instance should be under the same APPID. If this parameter is left empty, ID of the source instance will be used.",
			},

			"rename_restore": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Restore the databases listed in ReNameRestoreDatabase and rename them after restoration. If this parameter is left empty, all databases will be restored and renamed in the default format.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"old_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database name. If the OldName database does not exist, a failure will be returned.It can be left empty in offline migration tasks.",
						},
						"new_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "New database name. In offline migration, OldName will be used if NewName is left empty (OldName and NewName cannot be both empty). In database cloning, OldName and NewName must be both specified and cannot have the same value.",
						},
					},
				},
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Rollback type, 0-overwrite method; 1-rename method, default 1.",
			},

			"d_b_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The database that needs to be covered and rolled back is required only when the file is covered and rolled back.",
			},

			"group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Group ID of unarchived backup files grouped by backup task. This parameter is returned by the DescribeBackups API.",
			},
		},
	}
}

func resourceTencentCloudSqlserverRestoreInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restore_instance.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverRestoreInstanceUpdate(d, meta)
}

func resourceTencentCloudSqlserverRestoreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restore_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	restoreInstanceId := d.Id()

	restoreInstance, err := service.DescribeSqlserverRestoreInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if restoreInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverRestoreInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if restoreInstance.InstanceId != nil {
		_ = d.Set("instance_id", restoreInstance.InstanceId)
	}

	if restoreInstance.BackupId != nil {
		_ = d.Set("backup_id", restoreInstance.BackupId)
	}

	if restoreInstance.TargetInstanceId != nil {
		_ = d.Set("target_instance_id", restoreInstance.TargetInstanceId)
	}

	if restoreInstance.RenameRestore != nil {
		renameRestoreList := []interface{}{}
		for _, renameRestore := range restoreInstance.RenameRestore {
			renameRestoreMap := map[string]interface{}{}

			if restoreInstance.RenameRestore.OldName != nil {
				renameRestoreMap["old_name"] = restoreInstance.RenameRestore.OldName
			}

			if restoreInstance.RenameRestore.NewName != nil {
				renameRestoreMap["new_name"] = restoreInstance.RenameRestore.NewName
			}

			renameRestoreList = append(renameRestoreList, renameRestoreMap)
		}

		_ = d.Set("rename_restore", renameRestoreList)

	}

	if restoreInstance.Type != nil {
		_ = d.Set("type", restoreInstance.Type)
	}

	if restoreInstance.DBList != nil {
		_ = d.Set("d_b_list", restoreInstance.DBList)
	}

	if restoreInstance.GroupId != nil {
		_ = d.Set("group_id", restoreInstance.GroupId)
	}

	return nil
}

func resourceTencentCloudSqlserverRestoreInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restore_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewRestoreInstanceRequest()

	restoreInstanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "backup_id", "target_instance_id", "rename_restore", "type", "d_b_list", "group_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().RestoreInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver restoreInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverRestoreInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverRestoreInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restore_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

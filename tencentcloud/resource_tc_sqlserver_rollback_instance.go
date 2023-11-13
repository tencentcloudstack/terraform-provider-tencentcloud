/*
Provides a resource to create a sqlserver rollback_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_rollback_instance" "rollback_instance" {
  instance_id = "mssql-i1z41iwd"
  type = 0
  time = ""
  d_bs =
  target_instance_id = ""
  rename_restore {
		old_name = ""
		new_name = ""

  }
}
```

Import

sqlserver rollback_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_rollback_instance.rollback_instance rollback_instance_id
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

func resourceTencentCloudSqlserverRollbackInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverRollbackInstanceCreate,
		Read:   resourceTencentCloudSqlserverRollbackInstanceRead,
		Update: resourceTencentCloudSqlserverRollbackInstanceUpdate,
		Delete: resourceTencentCloudSqlserverRollbackInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Rollback type. 0: the database rolled back overwrites the original database; 1: the database rolled back is renamed and does not overwrite the original database.",
			},

			"time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Target time point for rollback.",
			},

			"d_bs": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Database to be rolled back.",
			},

			"target_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of the target instance to which the backup is restored. The target instance should be under the same APPID. If this parameter is left empty, ID of the source instance will be used.",
			},

			"rename_restore": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Rename the databases listed in ReNameRestoreDatabase. This parameter takes effect only when Type = 1 which indicates that backup rollback supports renaming databases. If it is left empty, databases will be renamed in the default format and the DBs parameter specifies the databases to be restored.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"old_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database name. If the OldName database does not exist, a failure will be returned. It can be left empty in offline migration tasks.",
						},
						"new_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "New database name. In offline migration, OldName will be used if NewName is left empty (OldName and NewName cannot be both empty). In database cloning, OldName and NewName must be both specified and cannot have the same value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSqlserverRollbackInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_rollback_instance.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverRollbackInstanceUpdate(d, meta)
}

func resourceTencentCloudSqlserverRollbackInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_rollback_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	rollbackInstanceId := d.Id()

	rollbackInstance, err := service.DescribeSqlserverRollbackInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if rollbackInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverRollbackInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if rollbackInstance.InstanceId != nil {
		_ = d.Set("instance_id", rollbackInstance.InstanceId)
	}

	if rollbackInstance.Type != nil {
		_ = d.Set("type", rollbackInstance.Type)
	}

	if rollbackInstance.Time != nil {
		_ = d.Set("time", rollbackInstance.Time)
	}

	if rollbackInstance.DBs != nil {
		_ = d.Set("d_bs", rollbackInstance.DBs)
	}

	if rollbackInstance.TargetInstanceId != nil {
		_ = d.Set("target_instance_id", rollbackInstance.TargetInstanceId)
	}

	if rollbackInstance.RenameRestore != nil {
		renameRestoreList := []interface{}{}
		for _, renameRestore := range rollbackInstance.RenameRestore {
			renameRestoreMap := map[string]interface{}{}

			if rollbackInstance.RenameRestore.OldName != nil {
				renameRestoreMap["old_name"] = rollbackInstance.RenameRestore.OldName
			}

			if rollbackInstance.RenameRestore.NewName != nil {
				renameRestoreMap["new_name"] = rollbackInstance.RenameRestore.NewName
			}

			renameRestoreList = append(renameRestoreList, renameRestoreMap)
		}

		_ = d.Set("rename_restore", renameRestoreList)

	}

	return nil
}

func resourceTencentCloudSqlserverRollbackInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_rollback_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewRollbackInstanceRequest()

	rollbackInstanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "type", "time", "d_bs", "target_instance_id", "rename_restore"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().RollbackInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver rollbackInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverRollbackInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverRollbackInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_rollback_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

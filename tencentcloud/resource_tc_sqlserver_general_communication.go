/*
Provides a resource to create a sqlserver general_communication

Example Usage

```hcl
resource "tencentcloud_sqlserver_general_communication" "general_communication" {
  instance_id = "Instance ID in the format of mssql-j8kv137v"
  rename_restore {
		old_name = "old_db_name"
		new_name = "new_db_name"

  }
}
```

Import

sqlserver general_communication can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_communication.general_communication general_communication_id
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

func resourceTencentCloudSqlserverGeneralCommunication() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralCommunicationCreate,
		Read:   resourceTencentCloudSqlserverGeneralCommunicationRead,
		Update: resourceTencentCloudSqlserverGeneralCommunicationUpdate,
		Delete: resourceTencentCloudSqlserverGeneralCommunicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"rename_restore": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Clone and rename the databases specified in ReNameRestoreDatabase. Please note that the clones must be renamed.",
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

func resourceTencentCloudSqlserverGeneralCommunicationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_communication.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = sqlserver.NewCloneDBRequest()
		response   = sqlserver.NewCloneDBResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rename_restore"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			renameRestoreDatabase := sqlserver.RenameRestoreDatabase{}
			if v, ok := dMap["old_name"]; ok {
				renameRestoreDatabase.OldName = helper.String(v.(string))
			}
			if v, ok := dMap["new_name"]; ok {
				renameRestoreDatabase.NewName = helper.String(v.(string))
			}
			request.RenameRestore = append(request.RenameRestore, &renameRestoreDatabase)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CloneDB(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver generalCommunication failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudSqlserverGeneralCommunicationRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCommunicationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_communication.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	generalCommunicationId := d.Id()

	generalCommunication, err := service.DescribeSqlserverGeneralCommunicationById(ctx, instanceId)
	if err != nil {
		return err
	}

	if generalCommunication == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverGeneralCommunication` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if generalCommunication.InstanceId != nil {
		_ = d.Set("instance_id", generalCommunication.InstanceId)
	}

	if generalCommunication.RenameRestore != nil {
		renameRestoreList := []interface{}{}
		for _, renameRestore := range generalCommunication.RenameRestore {
			renameRestoreMap := map[string]interface{}{}

			if generalCommunication.RenameRestore.OldName != nil {
				renameRestoreMap["old_name"] = generalCommunication.RenameRestore.OldName
			}

			if generalCommunication.RenameRestore.NewName != nil {
				renameRestoreMap["new_name"] = generalCommunication.RenameRestore.NewName
			}

			renameRestoreList = append(renameRestoreList, renameRestoreMap)
		}

		_ = d.Set("rename_restore", renameRestoreList)

	}

	return nil
}

func resourceTencentCloudSqlserverGeneralCommunicationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_communication.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyDBNameRequest()

	generalCommunicationId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "rename_restore"}

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyDBName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver generalCommunication failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverGeneralCommunicationRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCommunicationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_general_communication.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	generalCommunicationId := d.Id()

	if err := service.DeleteSqlserverGeneralCommunicationById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}

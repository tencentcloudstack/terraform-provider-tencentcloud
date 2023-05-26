/*
Provides a resource to create a sqlserver rollback_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_rollback_instance" "rollback_instance" {
  instance_id = "mssql-qelbzgwf"
  time = "%s"
  rename_restore {
    old_name = "keep_pubsub_db2"
	new_name = "rollback_pubsub_db3"
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
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
			"time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Target time point for rollback.",
			},
			"rename_restore": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Rename the databases listed in ReNameRestoreDatabase. This parameter takes effect only when Type = 1 which indicates that backup rollback supports renaming databases. If it is left empty, databases will be renamed in the default format and the DBs parameter specifies the databases to be restored.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"old_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name. If the OldName database does not exist, a failure will be returned. It can be left empty in offline migration tasks.",
						},
						"new_name": {
							Type:        schema.TypeString,
							Required:    true,
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

	var (
		instanceId string
		tmpTime    string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("time"); ok {
		tmpTime = v.(string)
	}

	oldNameList := make([]string, 0)
	newNameList := make([]string, 0)
	if v, ok := d.GetOk("rename_restore"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			if v, ok := dMap["old_name"]; ok {
				oldNameList = append(oldNameList, v.(string))
			}
			if v, ok := dMap["new_name"]; ok {
				newNameList = append(newNameList, v.(string))
			}
		}
	}

	oldNameListStr := strings.Join(oldNameList, COMMA_SP)
	newNameListStr := strings.Join(newNameList, COMMA_SP)

	d.SetId(strings.Join([]string{instanceId, tmpTime, oldNameListStr, newNameListStr}, FILED_SP))

	return resourceTencentCloudSqlserverRollbackInstanceUpdate(d, meta)
}

func resourceTencentCloudSqlserverRollbackInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_rollback_instance.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	tmpTime := idSplit[1]
	oldNameListStr := idSplit[2]
	newNameListStr := idSplit[3]
	oldNameList := strings.Split(oldNameListStr, COMMA_SP)
	newNameList := strings.Split(newNameListStr, COMMA_SP)

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

	_ = d.Set("time", tmpTime)

	renameRestoreList := []interface{}{}
	for i := 0; i < len(oldNameList); i++ {
		renameRestoreMap := map[string]interface{}{}
		renameRestoreMap["old_name"] = oldNameList[i]
		renameRestoreMap["new_name"] = newNameList[i]
		renameRestoreList = append(renameRestoreList, renameRestoreMap)
	}
	_ = d.Set("rename_restore", renameRestoreList)

	return nil
}

func resourceTencentCloudSqlserverRollbackInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_rollback_instance.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		request = sqlserver.NewRollbackInstanceRequest()
		flowId  uint64
		tmpType uint64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]

	request.InstanceId = &instanceId
	tmpType = 1
	request.Type = &tmpType

	if v, ok := d.GetOk("time"); ok {
		request.Time = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rename_restore"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			parameter := sqlserver.RenameRestoreDatabase{}
			if v, ok := dMap["old_name"]; ok {
				parameter.OldName = helper.String(v.(string))
			}
			if v, ok := dMap["new_name"]; ok {
				parameter.NewName = helper.String(v.(string))
			}
			request.RenameRestore = append(request.RenameRestore, &parameter)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().RollbackInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver rollbackInstance failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCloneStatusByFlowId(ctx, int64(flowId))
		if e != nil {
			return retryError(e)
		}

		if result == nil {
			e = fmt.Errorf("sqlserver rollbackInstance instanceId %s flowId %d not exists", instanceId, flowId)
			return resource.NonRetryableError(e)
		}

		if *result.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("sqlserver rollbackInstance task status is running"))
		}

		if *result.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		}

		if *result.Status == SQLSERVER_TASK_FAIL {
			return resource.NonRetryableError(fmt.Errorf("sqlserver rollbackInstance task status is failed"))
		}

		e = fmt.Errorf("sqlserver rollbackInstance task status is %v, we won't wait for it finish", *result.Status)
		return resource.NonRetryableError(e)
	})

	if err != nil {
		log.Printf("[CRITAL]%s sqlserver rollbackInstance task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudSqlserverRollbackInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverRollbackInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_rollback_instance.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId            = getLogId(contextNil)
		ctx              = context.WithValue(context.TODO(), logIdKey, logId)
		sqlserverService = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	newNameListStr := idSplit[3]
	newNameList := strings.Split(newNameListStr, COMMA_SP)

	if len(newNameList) == 0 {
		return nil
	}

	tmpNames := make([]*string, len(newNameList))
	for k, v := range newNameList {
		tmpNames[k] = &v
	}

	err := sqlserverService.DeleteSqlserverDB(ctx, instanceId, tmpNames)
	return err
}

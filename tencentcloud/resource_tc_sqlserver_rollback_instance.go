/*
Provides a resource to create a sqlserver rollback_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_rollback_instance" "rollback_instance" {
  instance_id = "mssql-qelbzgwf"
  type = 1
  time = "2023-05-25 19:14:30"
  dbs = ["keep_pubsub_db2"]
  rename_restore {
    old_name = "keep_pubsub_db2"
	new_name = "rollback_pubsub_db2"
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
	"strconv"
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
			"dbs": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Database to be rolled back.",
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

	var (
		instanceId string
		tmpTime    string
		tmpType    string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		tmpType = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("time"); ok {
		tmpTime = v.(string)
	}

	dbList := make([]string, 0)
	if v, ok := d.GetOk("dbs"); ok {
		for _, item := range v.(*schema.Set).List() {
			dbList = append(dbList, item.(string))
		}
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

	dbListStr := strings.Join(dbList, COMMA_SP)
	oldNameListStr := strings.Join(oldNameList, COMMA_SP)
	newNameListStr := strings.Join(newNameList, COMMA_SP)

	d.SetId(strings.Join([]string{instanceId, tmpType, tmpTime, dbListStr, oldNameListStr, newNameListStr}, FILED_SP))

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
	if len(idSplit) != 6 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	tmpType := idSplit[1]
	tmpTime := idSplit[2]
	dbListStr := idSplit[3]
	oldNameListStr := idSplit[4]
	newNameListStr := idSplit[5]
	dbList := strings.Split(dbListStr, COMMA_SP)
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

	Type, _ := strconv.Atoi(tmpType)
	_ = d.Set("type", Type)
	_ = d.Set("time", tmpTime)
	if len(dbList) != 0 {
		_ = d.Set("dbs", dbList)
	}
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
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 6 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]

	request.InstanceId = &instanceId
	v, _ := d.GetOk("type")
	if v == 0 || v == 1 {
		request.Type = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("dbs"); ok {
		for _, item := range v.(*schema.Set).List() {
			request.DBs = append(request.DBs, helper.String(item.(string)))
		}
	}

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

	return nil
}

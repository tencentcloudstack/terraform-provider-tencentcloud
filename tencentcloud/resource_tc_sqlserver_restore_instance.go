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
			"rename_restore": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Restore the databases listed in ReNameRestoreDatabase and rename them after restoration. If this parameter is left empty, all databases will be restored and renamed in the default format.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"old_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name. If the OldName database does not exist, a failure will be returned.It can be left empty in offline migration tasks.",
						},
						"new_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "New database name. In offline migration, OldName will be used if NewName is left empty (OldName and NewName cannot be both empty). In database cloning, OldName and NewName must be both specified and cannot have the same value.",
						},
					},
				},
			},
			"encryption": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "TDE encryption, `enable` encrypted, `disable` unencrypted.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "encryption, `enable` encrypted, `disable` unencrypted.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSqlserverRestoreInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restore_instance.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
		backupId   string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("backup_id"); ok {
		backupId = strconv.Itoa(v.(int))
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

	d.SetId(strings.Join([]string{instanceId, backupId, oldNameListStr, newNameListStr}, FILED_SP))

	return resourceTencentCloudSqlserverRestoreInstanceUpdate(d, meta)
}

func resourceTencentCloudSqlserverRestoreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restore_instance.read")()
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
	backupId := idSplit[1]
	oldNameListStr := idSplit[2]
	newNameListStr := idSplit[3]
	oldNameList := strings.Split(oldNameListStr, COMMA_SP)
	newNameList := strings.Split(newNameListStr, COMMA_SP)
	allNameList := append(oldNameList, newNameList...)
	restoreInstance, err := service.DescribeSqlserverRestoreInstanceById(ctx, instanceId, allNameList)
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

	tmpBackipId, _ := strconv.Atoi(backupId)
	_ = d.Set("backup_id", tmpBackipId)

	renameRestoreList := []interface{}{}
	for i := 0; i < len(oldNameList); i++ {
		renameRestoreMap := map[string]interface{}{}
		renameRestoreMap["old_name"] = oldNameList[i]
		renameRestoreMap["new_name"] = newNameList[i]
		renameRestoreList = append(renameRestoreList, renameRestoreMap)
	}
	_ = d.Set("rename_restore", renameRestoreList)

	if restoreInstance.DBDetails != nil {
		tmpList := make([]map[string]interface{}, 0)
		for _, item := range restoreInstance.DBDetails {
			dMap := map[string]interface{}{}
			if item.Name != nil {
				dMap["db_name"] = item.Name
			}
			if item.Encryption != nil {
				dMap["status"] = item.Encryption
			}
			tmpList = append(tmpList, dMap)
		}
		_ = d.Set("encryption", tmpList)
	}

	return nil
}

func resourceTencentCloudSqlserverRestoreInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restore_instance.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		request = sqlserver.NewRestoreInstanceRequest()
		flowId  int64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	backupId := idSplit[1]

	request.InstanceId = &instanceId
	tmpBackupId, _ := strconv.Atoi(backupId)
	request.BackupId = helper.IntInt64(tmpBackupId)

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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().RestoreInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver restoreInstance failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCloneStatusByFlowId(ctx, flowId)
		if e != nil {
			return retryError(e)
		}

		if result == nil {
			e = fmt.Errorf("sqlserver restoreInstance instanceId %s flowId %d not exists", instanceId, flowId)
			return resource.NonRetryableError(e)
		}

		if *result.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("restore sqlserver restoreInstance task status is running"))
		}

		if *result.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		}

		if *result.Status == SQLSERVER_TASK_FAIL {
			return resource.NonRetryableError(fmt.Errorf("restore sqlserver restoreInstance task status is failed"))
		}

		e = fmt.Errorf("restore sqlserver restoreInstance task status is %v, we won't wait for it finish", *result.Status)
		return resource.NonRetryableError(e)
	})

	if err != nil {
		log.Printf("[CRITAL]%s restore sqlserver restoreInstance task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudSqlserverRestoreInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverRestoreInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_restore_instance.delete")()
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
	for v := range newNameList {
		tmpNames[v] = &newNameList[v]
	}

	err := sqlserverService.DeleteSqlserverDB(ctx, instanceId, tmpNames)
	return err
}

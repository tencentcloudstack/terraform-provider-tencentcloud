package sqlserver

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverRollbackInstance() *schema.Resource {
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
				Description: "Rename the databases listed in ReNameRestoreDatabase.",
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
							Description: "New database name.",
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

func resourceTencentCloudSqlserverRollbackInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_rollback_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

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

	oldNameListStr := strings.Join(oldNameList, tccommon.COMMA_SP)
	newNameListStr := strings.Join(newNameList, tccommon.COMMA_SP)

	d.SetId(strings.Join([]string{instanceId, tmpTime, oldNameListStr, newNameListStr}, tccommon.FILED_SP))

	return resourceTencentCloudSqlserverRollbackInstanceUpdate(d, meta)
}

func resourceTencentCloudSqlserverRollbackInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_rollback_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	tmpTime := idSplit[1]
	oldNameListStr := idSplit[2]
	newNameListStr := idSplit[3]
	oldNameList := strings.Split(oldNameListStr, tccommon.COMMA_SP)
	newNameList := strings.Split(newNameListStr, tccommon.COMMA_SP)
	allNameList := append(oldNameList, newNameList...)
	rollbackInstance, err := service.DescribeSqlserverRollbackInstanceById(ctx, instanceId, allNameList)
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

	if rollbackInstance.DBDetails != nil {
		tmpList := make([]map[string]interface{}, 0)
		for _, item := range rollbackInstance.DBDetails {
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

func resourceTencentCloudSqlserverRollbackInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_rollback_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request = sqlserver.NewRollbackInstanceRequest()
		flowId  uint64
		tmpType uint64
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().RollbackInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	err = resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCloneStatusByFlowId(ctx, int64(flowId))
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_rollback_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		sqlserverService = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	newNameListStr := idSplit[3]
	newNameList := strings.Split(newNameListStr, tccommon.COMMA_SP)

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

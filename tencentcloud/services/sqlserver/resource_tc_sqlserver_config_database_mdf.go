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
)

func ResourceTencentCloudSqlserverConfigDatabaseMdf() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigDatabaseMdfCreate,
		Read:   resourceTencentCloudSqlserverConfigDatabaseMdfRead,
		Update: resourceTencentCloudSqlserverConfigDatabaseMdfUpdate,
		Delete: resourceTencentCloudSqlserverConfigDatabaseMdfDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Array of database names.",
			},
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigDatabaseMdfCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_database_mdf.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
		dbName     string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("db_name"); ok {
		dbName = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, dbName}, tccommon.FILED_SP))

	return resourceTencentCloudSqlserverConfigDatabaseMdfUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigDatabaseMdfRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_database_mdf.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		Name    string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	dbName := idSplit[1]

	configDatabaseMdf, err := service.DescribeSqlserverConfigDatabaseMdfById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configDatabaseMdf == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigDatabaseMdf` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	for _, i := range configDatabaseMdf {
		if *i.Name == dbName {
			Name = *i.Name
			break
		}
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("db_name", Name)

	return nil
}

func resourceTencentCloudSqlserverConfigDatabaseMdfUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_database_mdf.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = sqlserver.NewModifyDatabaseMdfRequest()
		flowRequest = sqlserver.NewDescribeFlowStatusRequest()
		flowId      int64
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	dbName := idSplit[1]

	request.InstanceId = &instanceId
	dbNames := make([]*string, 0)
	dbNames = append(dbNames, &dbName)
	request.DBNames = dbNames

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyDatabaseMdf(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver configDatabaseMdf not exists")
			return resource.NonRetryableError(e)
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configDatabaseMdf failed, reason:%+v", logId, err)
		return err
	}

	flowRequest.FlowId = &flowId
	err = resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().DescribeFlowStatus(flowRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if *result.Response.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		} else if *result.Response.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("sqlserver configDatabaseMdf status is running"))
		} else if *result.Response.Status == int64(SQLSERVER_TASK_FAIL) {
			return resource.NonRetryableError(fmt.Errorf("sqlserver configDatabaseMdf status is fail"))
		} else {
			e = fmt.Errorf("sqlserver configDatabaseMdf status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver configDatabaseMdf failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigDatabaseMdfRead(d, meta)
}

func resourceTencentCloudSqlserverConfigDatabaseMdfDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_database_mdf.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

package sqlserver

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
)

func ResourceTencentCloudSqlserverConfigTerminateDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigTerminateDBInstanceCreate,
		Read:   resourceTencentCloudSqlserverConfigTerminateDBInstanceRead,
		Update: resourceTencentCloudSqlserverConfigTerminateDBInstanceUpdate,
		Delete: resourceTencentCloudSqlserverConfigTerminateDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigTerminateDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_terminate_db_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverConfigTerminateDBInstanceUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigTerminateDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_terminate_db_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	configTerminateDBInstance, err := service.DescribeSqlserverConfigTerminateDBInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configTerminateDBInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigTerminateDBInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configTerminateDBInstance.InstanceId != nil {
		_ = d.Set("instance_id", configTerminateDBInstance.InstanceId)
	}

	return nil
}

func resourceTencentCloudSqlserverConfigTerminateDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_terminate_db_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = sqlserver.NewTerminateDBInstanceRequest()
		instanceId = d.Id()
	)

	request.InstanceIdSet = []*string{&instanceId}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().TerminateDBInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configTerminateDBInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigTerminateDBInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverConfigTerminateDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_terminate_db_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

package sqlserver

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
)

func ResourceTencentCloudSqlserverRenewPostpaidDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverRenewPostpaidDBInstanceCreate,
		Read:   resourceTencentCloudSqlserverRenewPostpaidDBInstanceRead,
		Update: resourceTencentCloudSqlserverRenewPostpaidDBInstanceUpdate,
		Delete: resourceTencentCloudSqlserverRenewPostpaidDBInstanceDelete,
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

func resourceTencentCloudSqlserverRenewPostpaidDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_renew_postpaid_db_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverRenewPostpaidDBInstanceUpdate(d, meta)
}

func resourceTencentCloudSqlserverRenewPostpaidDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_renew_postpaid_db_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	renewPostpaidDBInstance, err := service.DescribeSqlserverRenewPostpaidDBInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if renewPostpaidDBInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverRenewPostpaidDBInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if renewPostpaidDBInstance.InstanceId != nil {
		_ = d.Set("instance_id", renewPostpaidDBInstance.InstanceId)
	}

	return nil
}

func resourceTencentCloudSqlserverRenewPostpaidDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_renew_postpaid_db_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = sqlserver.NewRenewPostpaidDBInstanceRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().RenewPostpaidDBInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver renewPostpaidDBInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverRenewPostpaidDBInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverRenewPostpaidDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_renew_postpaid_db_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

package cdwch

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clickhouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClickhouseBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClickhouseBackupCreate,
		Read:   resourceTencentCloudClickhouseBackupRead,
		Update: resourceTencentCloudClickhouseBackupUpdate,
		Delete: resourceTencentCloudClickhouseBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"cos_bucket_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "COS bucket name.",
			},
		},
	}
}

func resourceTencentCloudClickhouseBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clickhouse_backup.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = clickhouse.NewOpenBackUpRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	request.OperationType = helper.String("open")

	if v, ok := d.GetOk("cos_bucket_name"); ok {
		request.CosBucketName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdwchClient().OpenBackUp(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create clickhouse backup failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudClickhouseBackupRead(d, meta)
}

func resourceTencentCloudClickhouseBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clickhouse_backup.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CdwchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	backup, err := service.DescribeBackUpScheduleById(ctx, instanceId)
	if err != nil {
		return err
	}

	if backup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClickhouseBackup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if backup.MetaStrategy != nil && backup.MetaStrategy.CosBucketName != nil {
		_ = d.Set("cos_bucket_name", backup.MetaStrategy.CosBucketName)
	}

	return nil
}

func resourceTencentCloudClickhouseBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clickhouse_backup.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"instance_id", "operation_type", "cos_bucket_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudClickhouseBackupRead(d, meta)
}

func resourceTencentCloudClickhouseBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clickhouse_backup.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	instanceId := d.Id()
	request := clickhouse.NewOpenBackUpRequest()
	request.InstanceId = helper.String(instanceId)
	request.OperationType = helper.String("close")
	request.CosBucketName = helper.String(d.Get("cos_bucket_name").(string))

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdwchClient().OpenBackUp(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create clickhouse backup failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

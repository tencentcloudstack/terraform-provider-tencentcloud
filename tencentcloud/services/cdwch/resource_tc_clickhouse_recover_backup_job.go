package cdwch

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clickhouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClickhouseRecoverBackupJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClickhouseRecoverBackupJobCreate,
		Read:   resourceTencentCloudClickhouseRecoverBackupJobRead,
		Delete: resourceTencentCloudClickhouseRecoverBackupJobDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"back_up_job_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Back up job id.",
			},
		},
	}
}

func resourceTencentCloudClickhouseRecoverBackupJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clickhouse_recover_backup_job.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request     = clickhouse.NewRecoverBackUpJobRequest()
		instanceId  string
		backUpJobId int
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, _ := d.GetOk("back_up_job_id"); v != nil {
		backUpJobId = v.(int)
		request.BackUpJobId = helper.IntInt64(backUpJobId)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdwchClient().RecoverBackUpJob(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate clickhouse RecoverBackupJob failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + helper.IntToStr(backUpJobId))

	return resourceTencentCloudClickhouseRecoverBackupJobRead(d, meta)
}

func resourceTencentCloudClickhouseRecoverBackupJobRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clickhouse_recover_backup_job.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudClickhouseRecoverBackupJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clickhouse_recover_backup_job.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

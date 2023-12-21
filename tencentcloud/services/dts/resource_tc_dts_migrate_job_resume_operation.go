package dts

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDtsMigrateJobResumeOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsMigrateJobResumeOperationCreate,
		Read:   resourceTencentCloudDtsMigrateJobResumeOperationRead,
		Delete: resourceTencentCloudDtsMigrateJobResumeOperationDelete,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "job id.",
			},

			"resume_option": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "resume mode: 1.clearData-Clear target data; 2.overwrite-The task is executed in overwrite mode; 3.normal-No extra action. Note that clearData and overwrite are valid only for redis links, normal is valid only for non-Redis links.",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateJobResumeOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job_resume_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = dts.NewResumeMigrateJobRequest()
		jobId   string
	)
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
		request.JobId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resume_option"); ok {
		request.ResumeOption = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDtsClient().ResumeMigrateJob(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dts migrateJobResumeOperation failed, reason:%+v", logId, err)
		return err
	}
	d.SetId(jobId)

	service := DtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"running", "readyComplete"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DtsMigrateJobResumeOperationStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsMigrateJobResumeOperationRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobResumeOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job_resume_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDtsMigrateJobResumeOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job_resume_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

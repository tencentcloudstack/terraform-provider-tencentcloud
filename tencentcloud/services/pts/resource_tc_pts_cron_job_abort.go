package pts

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPtsCronJobAbort() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsCronJobAbortCreate,
		Read:   resourceTencentCloudPtsCronJobAbortRead,
		Delete: resourceTencentCloudPtsCronJobAbortDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},

			"cron_job_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cron job ID.",
			},
		},
	}
}

func resourceTencentCloudPtsCronJobAbortCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_cron_job_abort.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = pts.NewAbortCronJobsRequest()
		projectId string
		cronJobId string
	)
	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cron_job_id"); ok {
		cronJobId = v.(string)
		request.CronJobIds = append(request.CronJobIds, helper.String(cronJobId))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePtsClient().AbortCronJobs(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate pts cronJobAbort failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(projectId + tccommon.FILED_SP + cronJobId)

	return resourceTencentCloudPtsCronJobAbortRead(d, meta)
}

func resourceTencentCloudPtsCronJobAbortRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_cron_job_abort.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPtsCronJobAbortDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_cron_job_abort.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

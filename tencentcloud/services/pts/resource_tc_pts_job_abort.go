package pts

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPtsJobAbort() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsJobAbortCreate,
		Read:   resourceTencentCloudPtsJobAbortRead,
		Delete: resourceTencentCloudPtsJobAbortDelete,

		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Job ID.",
			},

			"project_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},

			"scenario_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Scenario ID.",
			},

			"abort_reason": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The reason for aborting the job.",
			},
		},
	}
}

func resourceTencentCloudPtsJobAbortCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_job_abort.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = pts.NewAbortJobRequest()
		projectId  string
		scenarioId string
		jobId      string
	)
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
		request.JobId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scenario_id"); ok {
		scenarioId = v.(string)
		request.ScenarioId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("abort_reason"); v != nil {
		request.AbortReason = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePtsClient().AbortJob(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate pts jobAbort failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(projectId + tccommon.FILED_SP + scenarioId + tccommon.FILED_SP + jobId)

	return resourceTencentCloudPtsJobAbortRead(d, meta)
}

func resourceTencentCloudPtsJobAbortRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_job_abort.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPtsJobAbortDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_pts_job_abort.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

package oceanus

import (
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOceanusStopJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusStopJobCreate,
		Read:   resourceTencentCloudOceanusStopJobRead,
		Delete: resourceTencentCloudOceanusStopJobDelete,

		Schema: map[string]*schema.Schema{
			"stop_job_descriptions": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "The description information for batch job stop.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Job Id.",
						},
						"stop_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Stop type,1 stopped 2 paused.",
						},
					},
				},
			},
			"work_space_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
			},
		},
	}
}

func resourceTencentCloudOceanusStopJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_stop_job.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = oceanus.NewStopJobsRequest()
		jobIds  []string
	)

	if v, ok := d.GetOk("stop_job_descriptions"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			stopJobDescription := oceanus.StopJobDescription{}
			if v, ok := dMap["job_id"]; ok {
				stopJobDescription.JobId = helper.String(v.(string))
				jobIds = append(jobIds, v.(string))
			}

			if v, ok := dMap["stop_type"]; ok {
				stopJobDescription.StopType = helper.IntInt64(v.(int))
			}

			request.StopJobDescriptions = append(request.StopJobDescriptions, &stopJobDescription)
		}
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		request.WorkSpaceId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOceanusClient().StopJobs(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate oceanus stopJob failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join(jobIds, tccommon.FILED_SP))

	return resourceTencentCloudOceanusStopJobRead(d, meta)
}

func resourceTencentCloudOceanusStopJobRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_stop_job.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudOceanusStopJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_stop_job.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

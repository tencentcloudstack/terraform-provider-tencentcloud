package tencentcloud

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudOceanusStopJob() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_oceanus_stop_job.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOceanusClient().StopJobs(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate oceanus stopJob failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join(jobIds, FILED_SP))

	return resourceTencentCloudOceanusStopJobRead(d, meta)
}

func resourceTencentCloudOceanusStopJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_stop_job.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudOceanusStopJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_stop_job.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

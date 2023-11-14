/*
Provides a resource to create a pts operate_job

Example Usage

```hcl
resource "tencentcloud_pts_operate_job" "operate_job" {
  job_id = ""
  project_id = ""
  scenario_id = ""
  abort_reason =
}
```

Import

pts operate_job can be imported using the id, e.g.

```
terraform import tencentcloud_pts_operate_job.operate_job operate_job_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudPtsOperateJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsOperateJobCreate,
		Read:   resourceTencentCloudPtsOperateJobRead,
		Delete: resourceTencentCloudPtsOperateJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

func resourceTencentCloudPtsOperateJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_operate_job.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = pts.NewAbortJobRequest()
		response = pts.NewAbortJobResponse()
		jobId    string
	)
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
		request.JobId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scenario_id"); ok {
		request.ScenarioId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("abort_reason"); v != nil {
		request.AbortReason = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().AbortJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate pts operateJob failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobId
	d.SetId(jobId)

	return resourceTencentCloudPtsOperateJobRead(d, meta)
}

func resourceTencentCloudPtsOperateJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_operate_job.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPtsOperateJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_operate_job.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

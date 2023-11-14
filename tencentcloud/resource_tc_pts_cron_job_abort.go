/*
Provides a resource to create a pts cron_job_abort

Example Usage

```hcl
resource "tencentcloud_pts_cron_job_abort" "cron_job_abort" {
  project_id = ""
  cron_job_ids =
}
```

Import

pts cron_job_abort can be imported using the id, e.g.

```
terraform import tencentcloud_pts_cron_job_abort.cron_job_abort cron_job_abort_id
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

func resourceTencentCloudPtsCronJobAbort() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsCronJobAbortCreate,
		Read:   resourceTencentCloudPtsCronJobAbortRead,
		Delete: resourceTencentCloudPtsCronJobAbortDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},

			"cron_job_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Cron job ID.",
			},
		},
	}
}

func resourceTencentCloudPtsCronJobAbortCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_cron_job_abort.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = pts.NewAbortCronJobsRequest()
		response = pts.NewAbortCronJobsResponse()
		jobId    string
	)
	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cron_job_ids"); ok {
		cronJobIdsSet := v.(*schema.Set).List()
		for i := range cronJobIdsSet {
			cronJobIds := cronJobIdsSet[i].(string)
			request.CronJobIds = append(request.CronJobIds, &cronJobIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().AbortCronJobs(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate pts cronJobAbort failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobId
	d.SetId(jobId)

	return resourceTencentCloudPtsCronJobAbortRead(d, meta)
}

func resourceTencentCloudPtsCronJobAbortRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_cron_job_abort.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPtsCronJobAbortDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_cron_job_abort.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

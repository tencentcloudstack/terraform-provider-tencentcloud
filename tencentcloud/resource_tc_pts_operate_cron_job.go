/*
Provides a resource to create a pts operate_cron_job

Example Usage

```hcl
resource "tencentcloud_pts_operate_cron_job" "operate_cron_job" {
  project_id = "project-abc"
  cron_job_ids =
}
```

Import

pts operate_cron_job can be imported using the id, e.g.

```
terraform import tencentcloud_pts_operate_cron_job.operate_cron_job operate_cron_job_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudPtsOperateCronJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsOperateCronJobCreate,
		Read:   resourceTencentCloudPtsOperateCronJobRead,
		Delete: resourceTencentCloudPtsOperateCronJobDelete,
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
				Description: "Cron job ID list.",
			},
		},
	}
}

func resourceTencentCloudPtsOperateCronJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_operate_cron_job.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = pts.NewRestartCronJobsRequest()
		response   = pts.NewRestartCronJobsResponse()
		projectId  string
		cronJobIds string
	)
	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().RestartCronJobs(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate pts operateCronJob failed, reason:%+v", logId, err)
		return err
	}

	projectId = *response.Response.ProjectId
	d.SetId(strings.Join([]string{projectId, cronJobIds}, FILED_SP))

	return resourceTencentCloudPtsOperateCronJobRead(d, meta)
}

func resourceTencentCloudPtsOperateCronJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_operate_cron_job.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPtsOperateCronJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_operate_cron_job.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

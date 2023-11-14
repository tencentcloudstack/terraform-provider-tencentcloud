/*
Provides a resource to create a pts cron_job

Example Usage

```hcl
resource "tencentcloud_pts_cron_job" "cron_job" {
  name = &lt;nil&gt;
  project_id = &lt;nil&gt;
  scenario_id = &lt;nil&gt;
  scenario_name = &lt;nil&gt;
  frequency_type = &lt;nil&gt;
  cron_expression = &lt;nil&gt;
  job_owner = &lt;nil&gt;
  end_time = &lt;nil&gt;
  notice_id = &lt;nil&gt;
  note = &lt;nil&gt;
              }
```

Import

pts cron_job can be imported using the id, e.g.

```
terraform import tencentcloud_pts_cron_job.cron_job cron_job_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudPtsCronJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsCronJobCreate,
		Read:   resourceTencentCloudPtsCronJobRead,
		Update: resourceTencentCloudPtsCronJobUpdate,
		Delete: resourceTencentCloudPtsCronJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cron Job Name.",
			},

			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project Id.",
			},

			"scenario_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Scenario Id.",
			},

			"scenario_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Scenario Name.",
			},

			"frequency_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Execution frequency type, `1`: execute only once; `2`: daily granularity; `3`: weekly granularity; `4`: advanced.",
			},

			"cron_expression": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cron expression.",
			},

			"job_owner": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job Owner.",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "End Time; type: Timestamp ISO8601.",
			},

			"notice_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Notice ID.",
			},

			"note": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Note.",
			},

			"abort_reason": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Reason for suspension.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Scheduled task status.",
			},

			"created_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creation time; type: Timestamp ISO8601.",
			},

			"updated_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Update time; type: Timestamp ISO8601.",
			},

			"app_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "App ID.",
			},

			"uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "User ID.",
			},

			"sub_account_uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Sub-user ID.",
			},
		},
	}
}

func resourceTencentCloudPtsCronJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_cron_job.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = pts.NewCreateCronJobRequest()
		response  = pts.NewCreateCronJobResponse()
		cronJobId string
		projectId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scenario_id"); ok {
		request.ScenarioId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scenario_name"); ok {
		request.ScenarioName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("frequency_type"); ok {
		request.FrequencyType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("cron_expression"); ok {
		request.CronExpression = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_owner"); ok {
		request.JobOwner = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("notice_id"); ok {
		request.NoticeId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("note"); ok {
		request.Note = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().CreateCronJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create pts cronJob failed, reason:%+v", logId, err)
		return err
	}

	cronJobId = *response.Response.CronJobId
	d.SetId(strings.Join([]string{cronJobId, projectId}, FILED_SP))

	return resourceTencentCloudPtsCronJobRead(d, meta)
}

func resourceTencentCloudPtsCronJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_cron_job.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	cronJobId := idSplit[0]
	projectId := idSplit[1]

	cronJob, err := service.DescribePtsCronJobById(ctx, cronJobId, projectId)
	if err != nil {
		return err
	}

	if cronJob == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PtsCronJob` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cronJob.Name != nil {
		_ = d.Set("name", cronJob.Name)
	}

	if cronJob.ProjectId != nil {
		_ = d.Set("project_id", cronJob.ProjectId)
	}

	if cronJob.ScenarioId != nil {
		_ = d.Set("scenario_id", cronJob.ScenarioId)
	}

	if cronJob.ScenarioName != nil {
		_ = d.Set("scenario_name", cronJob.ScenarioName)
	}

	if cronJob.FrequencyType != nil {
		_ = d.Set("frequency_type", cronJob.FrequencyType)
	}

	if cronJob.CronExpression != nil {
		_ = d.Set("cron_expression", cronJob.CronExpression)
	}

	if cronJob.JobOwner != nil {
		_ = d.Set("job_owner", cronJob.JobOwner)
	}

	if cronJob.EndTime != nil {
		_ = d.Set("end_time", cronJob.EndTime)
	}

	if cronJob.NoticeId != nil {
		_ = d.Set("notice_id", cronJob.NoticeId)
	}

	if cronJob.Note != nil {
		_ = d.Set("note", cronJob.Note)
	}

	if cronJob.AbortReason != nil {
		_ = d.Set("abort_reason", cronJob.AbortReason)
	}

	if cronJob.Status != nil {
		_ = d.Set("status", cronJob.Status)
	}

	if cronJob.CreatedAt != nil {
		_ = d.Set("created_at", cronJob.CreatedAt)
	}

	if cronJob.UpdatedAt != nil {
		_ = d.Set("updated_at", cronJob.UpdatedAt)
	}

	if cronJob.AppId != nil {
		_ = d.Set("app_id", cronJob.AppId)
	}

	if cronJob.Uin != nil {
		_ = d.Set("uin", cronJob.Uin)
	}

	if cronJob.SubAccountUin != nil {
		_ = d.Set("sub_account_uin", cronJob.SubAccountUin)
	}

	return nil
}

func resourceTencentCloudPtsCronJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_cron_job.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := pts.NewUpdateCronJobRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	cronJobId := idSplit[0]
	projectId := idSplit[1]

	request.CronJobId = &cronJobId
	request.ProjectId = &projectId

	immutableArgs := []string{"name", "project_id", "scenario_id", "scenario_name", "frequency_type", "cron_expression", "job_owner", "end_time", "notice_id", "note", "abort_reason", "status", "created_at", "updated_at", "app_id", "uin", "sub_account_uin"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOk("project_id"); ok {
			request.ProjectId = helper.String(v.(string))
		}
	}

	if d.HasChange("scenario_id") {
		if v, ok := d.GetOk("scenario_id"); ok {
			request.ScenarioId = helper.String(v.(string))
		}
	}

	if d.HasChange("scenario_name") {
		if v, ok := d.GetOk("scenario_name"); ok {
			request.ScenarioName = helper.String(v.(string))
		}
	}

	if d.HasChange("frequency_type") {
		if v, ok := d.GetOkExists("frequency_type"); ok {
			request.FrequencyType = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("cron_expression") {
		if v, ok := d.GetOk("cron_expression"); ok {
			request.CronExpression = helper.String(v.(string))
		}
	}

	if d.HasChange("job_owner") {
		if v, ok := d.GetOk("job_owner"); ok {
			request.JobOwner = helper.String(v.(string))
		}
	}

	if d.HasChange("end_time") {
		if v, ok := d.GetOk("end_time"); ok {
			request.EndTime = helper.String(v.(string))
		}
	}

	if d.HasChange("notice_id") {
		if v, ok := d.GetOk("notice_id"); ok {
			request.NoticeId = helper.String(v.(string))
		}
	}

	if d.HasChange("note") {
		if v, ok := d.GetOk("note"); ok {
			request.Note = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().UpdateCronJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update pts cronJob failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPtsCronJobRead(d, meta)
}

func resourceTencentCloudPtsCronJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_cron_job.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PtsService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	cronJobId := idSplit[0]
	projectId := idSplit[1]

	if err := service.DeletePtsCronJobById(ctx, cronJobId, projectId); err != nil {
		return err
	}

	return nil
}

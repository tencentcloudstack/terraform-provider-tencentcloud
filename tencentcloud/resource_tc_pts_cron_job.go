/*
Provides a resource to create a pts cron_job

Example Usage

```hcl
resource "tencentcloud_pts_cron_job" "cron_job" {
  name = "iac-cron_job-update"
  project_id = "project-7qkzxhea"
  scenario_id = "scenario-c22lqb1w"
  scenario_name = "pts-js(2022-11-10 21:53:53)"
  frequency_type = 2
  cron_expression = "* 1 * * *"
  job_owner = "userName"
  # end_time = ""
  notice_id = "notice-vp6i38jt"
  note = "desc"
}

```
Import

pts cron_job can be imported using the projectId#cronJobId, e.g.
```
$ terraform import tencentcloud_pts_cron_job.cron_job project-7qkzxhea#scenario-c22lqb1w
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPtsCronJob() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudPtsCronJobRead,
		Create: resourceTencentCloudPtsCronJobCreate,
		Update: resourceTencentCloudPtsCronJobUpdate,
		Delete: resourceTencentCloudPtsCronJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cron Job Name.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project Id.",
			},

			"scenario_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Scenario Id.",
			},

			"scenario_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Scenario Name.",
			},

			"frequency_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Execution frequency type, `1`: execute only once; `2`: daily granularity; `3`: weekly granularity; `4`: advanced.",
			},

			"cron_expression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cron expression, When setting cron_expression at that time, frequency_type must be greater than 1.",
			},

			"job_owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Job Owner.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "End Time; type: Timestamp ISO8601.",
			},

			"notice_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notice ID.",
			},

			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Note.",
			},

			"abort_reason": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Reason for suspension.",
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Scheduled task status.",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time; type: Timestamp ISO8601.",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time; type: Timestamp ISO8601.",
			},

			"app_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "App ID.",
			},

			"uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User ID.",
			},

			"cron_job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cron job ID.",
			},

			"sub_account_uin": {
				Type:        schema.TypeString,
				Computed:    true,
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
		response  *pts.CreateCronJobResponse
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

	if v, ok := d.GetOk("frequency_type"); ok {
		request.FrequencyType = helper.Int64(int64(v.(int)))
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create pts cronJob failed, reason:%+v", logId, err)
		return err
	}

	cronJobId = *response.Response.CronJobId

	d.SetId(projectId + FILED_SP + cronJobId)
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
	projectId := idSplit[0]
	cronJobId := idSplit[1]

	cronJob, err := service.DescribePtsCronJob(ctx, cronJobId, projectId)

	if err != nil {
		return err
	}

	if cronJob == nil {
		d.SetId("")
		return fmt.Errorf("resource `cronJob` %s does not exist", cronJobId)
	}

	_ = d.Set("cron_job_id", cronJobId)

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
	projectId := idSplit[0]
	cronJobId := idSplit[1]

	request.CronJobId = &cronJobId
	request.ProjectId = &projectId

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scenario_id"); ok {
		request.ScenarioId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scenario_name"); ok {
		request.ScenarioName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("frequency_type"); ok {
		request.FrequencyType = helper.Int64(int64(v.(int)))
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().UpdateCronJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create pts cronJob failed, reason:%+v", logId, err)
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
	projectId := idSplit[0]
	cronJobId := idSplit[1]

	if err := service.DeletePtsCronJobById(ctx, cronJobId, projectId); err != nil {
		return err
	}

	return nil
}

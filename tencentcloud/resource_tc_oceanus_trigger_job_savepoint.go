package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudOceanusTriggerJobSavepoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusTriggerJobSavepointCreate,
		Read:   resourceTencentCloudOceanusTriggerJobSavepointRead,
		Delete: resourceTencentCloudOceanusTriggerJobSavepointDelete,

		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Job SerialId.",
			},
			"description": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Savepoint description.",
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

func resourceTencentCloudOceanusTriggerJobSavepointCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_trigger_job_savepoint.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = oceanus.NewTriggerJobSavepointRequest()
		jobId   string
	)

	if v, ok := d.GetOk("job_id"); ok {
		request.JobId = helper.String(v.(string))
		jobId = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		request.WorkSpaceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOceanusClient().TriggerJobSavepoint(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate oceanus TriggerJobSavepoint failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(jobId)

	return resourceTencentCloudOceanusTriggerJobSavepointRead(d, meta)
}

func resourceTencentCloudOceanusTriggerJobSavepointRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_trigger_job_savepoint.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudOceanusTriggerJobSavepointDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_trigger_job_savepoint.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

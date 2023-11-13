/*
Provides a resource to create a dts sync_job_resume

Example Usage

```hcl
resource "tencentcloud_dts_sync_job_resume" "sync_job_resume" {
  job_id = "sync-werwfs23"
}
```

Import

dts sync_job_resume can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_job_resume.sync_job_resume sync_job_resume_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudDtsSyncJobResume() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsSyncJobResumeCreate,
		Read:   resourceTencentCloudDtsSyncJobResumeRead,
		Delete: resourceTencentCloudDtsSyncJobResumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Synchronization instance id (i.e. identifies a synchronization job).",
			},
		},
	}
}

func resourceTencentCloudDtsSyncJobResumeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_resume.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = dts.NewResumeSyncJobRequest()
		response = dts.NewResumeSyncJobResponse()
		jobId    string
	)
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
		request.JobId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().ResumeSyncJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dts syncJobResume failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobId
	d.SetId(jobId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Running"}, 2*readRetryTimeout, time.Second, service.DtsSyncJobResumeStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsSyncJobResumeRead(d, meta)
}

func resourceTencentCloudDtsSyncJobResumeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_resume.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDtsSyncJobResumeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_resume.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

/*
Provides a resource to create a dts migrate_job_resume

Example Usage

```hcl
resource "tencentcloud_dts_migrate_job_resume" "migrate_job_resume" {
  job_id = "dts-ekmhr27i"
  resume_option = "normal"
}
```

Import

dts migrate_job_resume can be imported using the id, e.g.

```
terraform import tencentcloud_dts_migrate_job_resume.migrate_job_resume migrate_job_resume_id
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

func resourceTencentCloudDtsMigrateJobResume() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsMigrateJobResumeCreate,
		Read:   resourceTencentCloudDtsMigrateJobResumeRead,
		Delete: resourceTencentCloudDtsMigrateJobResumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Job id.",
			},

			"resume_option": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Resume mode: clearData-Clear target data；overwrite-The task is executed in overwrite mode；normal-No extra action，for normal.Note that clearData and overwrite are valid only for redis links, normal is valid only for non-Redis links.",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateJobResumeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_resume.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = dts.NewResumeMigrateJobRequest()
		response = dts.NewResumeMigrateJobResponse()
		jobId    string
	)
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
		request.JobId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resume_option"); ok {
		request.ResumeOption = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().ResumeMigrateJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dts migrateJobResume failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobId
	d.SetId(jobId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running,readyComplete"}, 3*readRetryTimeout, time.Second, service.DtsMigrateJobResumeStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsMigrateJobResumeRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobResumeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_resume.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDtsMigrateJobResumeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_resume.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

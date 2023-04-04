/*
Provides a resource to create a dts sync_job_resume_operation

Example Usage

```hcl
resource "tencentcloud_dts_sync_job_resume_operation" "sync_job_resume_operation" {
  job_id = "sync-werwfs23"
}
```

*/
package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDtsSyncJobResumeOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsSyncJobResumeOperationCreate,
		Read:   resourceTencentCloudDtsSyncJobResumeOperationRead,
		Update: resourceTencentCloudDtsSyncJobResumeOperationUpdate,
		Delete: resourceTencentCloudDtsSyncJobResumeOperationDelete,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Synchronization instance id (i.e. identifies a synchronization job).",
			},
		},
	}
}

func resourceTencentCloudDtsSyncJobResumeOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_resume_operation.create")()
	defer inconsistentCheck(d, meta)()

	var jobId string

	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}
	d.SetId(jobId)

	return resourceTencentCloudDtsSyncJobResumeOperationUpdate(d, meta)
}

func resourceTencentCloudDtsSyncJobResumeOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_resume_operation.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	jobId := d.Id()

	syncJobResumeOperation, err := service.DescribeDtsSyncJobResumeOperationById(ctx, jobId)
	if err != nil {
		return err
	}

	if syncJobResumeOperation == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsSyncJobResumeOperation` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if len(syncJobResumeOperation.JobList) == 0 || *syncJobResumeOperation.TotalCount == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsSyncJobResumeOperation.JobList` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	return nil
}

func resourceTencentCloudDtsSyncJobResumeOperationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_resume_operation.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dts.NewResumeSyncJobRequest()

	jobId := d.Id()

	request.JobId = helper.String(jobId)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().ResumeSyncJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dts syncJobResumeOperation failed, reason:%+v", logId, err)
		return err
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Running"}, 2*readRetryTimeout, time.Second, service.DtsSyncJobResumeOperationStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsSyncJobResumeOperationRead(d, meta)
}

func resourceTencentCloudDtsSyncJobResumeOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_resume_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

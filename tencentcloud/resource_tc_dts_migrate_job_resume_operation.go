/*
Provides a resource to create a dts migrate_job_resume_operation

Example Usage

```hcl
resource "tencentcloud_dts_migrate_job_resume_operation" "migrate_job_resume_operation" {
  job_id = "dts-ekmhr27i"
  resume_option = "normal"
}
```

Import

dts migrate_job_resume_operation can be imported using the id, e.g.

```
terraform import tencentcloud_dts_migrate_job_resume_operation.migrate_job_resume_operation migrate_job_resume_operation_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDtsMigrateJobResumeOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsMigrateJobResumeOperationCreate,
		Read:   resourceTencentCloudDtsMigrateJobResumeOperationRead,
		Update: resourceTencentCloudDtsMigrateJobResumeOperationUpdate,
		Delete: resourceTencentCloudDtsMigrateJobResumeOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "job id.",
			},

			"resume_option": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "resume mode: 1.clearData-Clear target data; 2.overwrite-The task is executed in overwrite mode; 3.normal-No extra action. Note that clearData and overwrite are valid only for redis links, normal is valid only for non-Redis links.",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateJobResumeOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_resume_operation.create")()
	defer inconsistentCheck(d, meta)()

	var jobId string

	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}
	d.SetId(jobId)

	return resourceTencentCloudDtsMigrateJobResumeOperationRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobResumeOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_resume_operation.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	jobId := d.Id()

	migrateJobResumeOperation, err := service.DescribeDtsMigrateJobResumeOperationById(ctx, jobId)
	if err != nil {
		return err
	}

	if migrateJobResumeOperation == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsMigrateJobResumeOperation` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if migrateJobResumeOperation.JobId != nil {
		_ = d.Set("job_id", migrateJobResumeOperation.JobId)
	}

	// operation do not need to check the resume option
	// if migrateJobResumeOperation.ResumeOption != nil {
	// 	_ = d.Set("resume_option", migrateJobResumeOperation.ResumeOption)
	// }

	return nil
}

func resourceTencentCloudDtsMigrateJobResumeOperationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_resume_operation.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dts.NewResumeMigrateJobRequest()

	request.JobId = helper.String(d.Id())

	immutableArgs := []string{"job_id", "resume_option"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().ResumeMigrateJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dts migrateJobResumeOperation failed, reason:%+v", logId, err)
		return err
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running,readyComplete"}, 3*readRetryTimeout, time.Second, service.DtsMigrateJobResumeOperationStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsMigrateJobResumeOperationRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobResumeOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_resume_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

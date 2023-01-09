/*
Provides a resource to start a dts migrate_job

Example Usage

```hcl
resource "tencentcloud_dts_migrate_job_start_operation" "start"{
	job_id = tencentcloud_dts_migrate_job.job.id
}
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

func resourceTencentCloudDtsMigrateJobStartOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsMigrateJobStartOperationCreate,
		Read:   resourceTencentCloudDtsMigrateJobStartOperationRead,
		Delete: resourceTencentCloudDtsMigrateJobStartOperationDelete,
		Schema: map[string]*schema.Schema{
			"job_id": {
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job Id from `tencentcloud_dts_migrate_job`.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Migrate job status.",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateJobStartOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_start_operation.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		tcClient = meta.(*TencentCloudClient).apiV3Conn
		service  = DtsService{client: tcClient}
		jobId    string
	)

	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}

	log.Printf("[DEBUG]%s trying to Start Migrate job, jobId:[%s]", logId, jobId)
	object, err := service.DescribeDtsMigrateCheckById(ctx, jobId)
	if err != nil {
		log.Printf("[CRITAL]%s DescribeDtsMigrateCheckById failed, reason:%+v", logId, err)
		return err
	}

	result := *object.CheckFlag
	if result == "checkNotPass" {
		var errorLog []string
		for _, step := range object.StepInfo {
			if *step.StepStatus == "failed" {
				for i, item := range step.DetailCheckItems {
					errorLog = append(errorLog, fmt.Sprintf("[error %v] %s\n", i, *item.ErrorLog[0]))
				}
			}
		}
		return fmt.Errorf("The DTS migration check not passed. Please re-check the job. check_result:[%s], reason:[%s]", result, errorLog)
	}

	startMigrateJobRequest := dts.NewStartMigrateJobRequest()
	startMigrateJobRequest.JobId = helper.String(jobId)
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := tcClient.UseDtsClient().StartMigrateJob(startMigrateJobRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, startMigrateJobRequest.GetAction(), startMigrateJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s start dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	conf := BuildStateChangeConf([]string{}, []string{"running", "error"}, 3*readRetryTimeout, time.Second, service.DtsMigrateJobStateRefreshFunc(jobId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(jobId)
	return resourceTencentCloudDtsMigrateJobStartOperationRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobStartOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_start_operation.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	jobId := d.Id()

	migrateJob, err := service.DescribeDtsMigrateJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if migrateJob == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if migrateJob.JobId != nil {
		_ = d.Set("job_id", migrateJob.JobId)
	}

	if migrateJob.Status != nil {
		_ = d.Set("status", migrateJob.Status)
	}

	return nil
}

func resourceTencentCloudDtsMigrateJobStartOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_start_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

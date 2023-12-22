package dts

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDtsMigrateJobStartOperation() *schema.Resource {
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
		},
	}
}

func resourceTencentCloudDtsMigrateJobStartOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job_start_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
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
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := tcClient.UseDtsClient().StartMigrateJob(startMigrateJobRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, startMigrateJobRequest.GetAction(), startMigrateJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s start dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"running", "error"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DtsMigrateJobStateRefreshFunc(jobId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(jobId)
	return resourceTencentCloudDtsMigrateJobStartOperationRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobStartOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job_start_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDtsMigrateJobStartOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job_start_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

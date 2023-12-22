package dts

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDtsSyncJobIsolateOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsSyncJobIsolateOperationCreate,
		Read:   resourceTencentCloudDtsSyncJobIsolateOperationRead,
		Delete: resourceTencentCloudDtsSyncJobIsolateOperationDelete,
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

func resourceTencentCloudDtsSyncJobIsolateOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_job_isolate_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = dts.NewIsolateSyncJobRequest()
		jobId   string
	)
	if v, ok := d.GetOk("job_id"); ok {
		request.JobId = helper.String(v.(string))
		jobId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDtsClient().IsolateSyncJob(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dts syncJobIsolateOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(jobId)

	service := DtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"Isolated"}, 2*tccommon.ReadRetryTimeout, time.Second, service.DtsSyncJobTradeStateRefreshFunc(d.Id(), "Isolated", []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsSyncJobIsolateOperationRead(d, meta)
}

func resourceTencentCloudDtsSyncJobIsolateOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_job_isolate_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDtsSyncJobIsolateOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_job_isolate_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

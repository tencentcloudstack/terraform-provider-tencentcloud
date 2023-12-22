package dts

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDtsCompareTaskStopOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsCompareTaskStopOperationCreate,
		Read:   resourceTencentCloudDtsCompareTaskStopOperationRead,
		Delete: resourceTencentCloudDtsCompareTaskStopOperationDelete,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "job id.",
			},

			"compare_task_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Compare task id.",
			},
		},
	}
}

func resourceTencentCloudDtsCompareTaskStopOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_compare_task_stop_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request       = dts.NewStopCompareRequest()
		jobId         string
		compareTaskId string
	)
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
		request.JobId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("compare_task_id"); ok {
		compareTaskId = v.(string)
		request.CompareTaskId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_id"); ok {
		request.JobId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDtsClient().StopCompare(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dts compareTaskStopOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(jobId + tccommon.FILED_SP + compareTaskId)

	return resourceTencentCloudDtsCompareTaskStopOperationRead(d, meta)
}

func resourceTencentCloudDtsCompareTaskStopOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_compare_task_stop_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDtsCompareTaskStopOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_compare_task_stop_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

package dcdb

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDcdbCancelDcnJobOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbCancelDcnJobOperationCreate,
		Read:   resourceTencentCloudDcdbCancelDcnJobOperationRead,
		Delete: resourceTencentCloudDcdbCancelDcnJobOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudDcdbCancelDcnJobOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_cancel_dcn_job_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = dcdb.NewCancelDcnJobRequest()
		instanceId string
		flowId     *int64
		service    = DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().CancelDcnJob(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		flowId = result.Response.FlowId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb cancelDcnJobOperation failed, reason:%+v", logId, err)
		return err
	}

	// need to wait flow success
	// 0:success; 1:failed, 2:running
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"0"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DcdbDbInstanceStateRefreshFunc(flowId, []string{"1"}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(instanceId)

	return resourceTencentCloudDcdbCancelDcnJobOperationRead(d, meta)
}

func resourceTencentCloudDcdbCancelDcnJobOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_cancel_dcn_job_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbCancelDcnJobOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_cancel_dcn_job_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

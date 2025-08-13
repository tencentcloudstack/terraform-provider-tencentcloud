package dlc

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcSuspendResumeDataEngine() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcSuspendResumeDataEngineCreate,
		Read:   resourceTencentCloudDlcSuspendResumeDataEngineRead,
		Delete: resourceTencentCloudDlcSuspendResumeDataEngineDelete,
		Schema: map[string]*schema.Schema{
			"data_engine_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The name of a virtual cluster.",
			},

			"operate": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The operation type: `suspend` or `resume`.",
			},
		},
	}
}

func resourceTencentCloudDlcSuspendResumeDataEngineCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_suspend_resume_data_engine.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		request        = dlc.NewSuspendResumeDataEngineRequest()
		dataEngineName string
	)

	if v, ok := d.GetOk("data_engine_name"); ok {
		dataEngineName = v.(string)
		request.DataEngineName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("operate"); ok {
		request.Operate = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().SuspendResumeDataEngine(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate dlc suspendResumeDataEngine failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dataEngineName)
	return resourceTencentCloudDlcSuspendResumeDataEngineRead(d, meta)
}

func resourceTencentCloudDlcSuspendResumeDataEngineRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_suspend_resume_data_engine.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcSuspendResumeDataEngineDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_suspend_resume_data_engine.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

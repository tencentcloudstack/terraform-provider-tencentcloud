package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDlcSuspendResumeDataEngine() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcSuspendResumeDataEngineCreate,
		Read:   resourceTencentCloudDlcSuspendResumeDataEngineRead,
		Delete: resourceTencentCloudDlcSuspendResumeDataEngineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"data_engine_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Engine name.",
			},

			"operate": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Engine operate tye: suspend/resume.",
			},
		},
	}
}

func resourceTencentCloudDlcSuspendResumeDataEngineCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_suspend_resume_data_engine.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().SuspendResumeDataEngine(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_dlc_suspend_resume_data_engine.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcSuspendResumeDataEngineDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_suspend_resume_data_engine.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

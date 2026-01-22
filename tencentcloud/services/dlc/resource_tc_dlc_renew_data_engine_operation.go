package dlc

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcRenewDataEngineOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcRenewDataEngineCreate,
		Read:   resourceTencentCloudDlcRenewDataEngineRead,
		Delete: resourceTencentCloudDlcRenewDataEngineDelete,
		Schema: map[string]*schema.Schema{
			"data_engine_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CU queue name.",
			},

			"time_span": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Renewal period in months, which is at least one month.",
			},

			"pay_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Payment type. It is 1 by default and is prepaid.",
			},

			"time_unit": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unit. It is m by default, and only m can be filled in.",
			},

			"renew_flag": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Auto-renewal flag: 0 means the initial status, and there is no automatic renewal by default. If the user has the privilege to retain services with prepayment, there will be an automatic renewal. 1 means that there is an automatic renewal. 2 means that there is surely no automatic renewal. If it is not specified, the parameter is 0 by default.",
			},
		},
	}
}

func resourceTencentCloudDlcRenewDataEngineCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_renew_data_engine_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		request        = dlc.NewRenewDataEngineRequest()
		dataEngineName string
	)

	if v, ok := d.GetOk("data_engine_name"); ok {
		dataEngineName = v.(string)
		request.DataEngineName = helper.String(v.(string))
	}

	if v, _ := d.GetOkExists("time_span"); v != nil {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOkExists("pay_mode"); v != nil {
		request.PayMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("time_unit"); ok {
		request.TimeUnit = helper.String(v.(string))
	}

	if v, _ := d.GetOkExists("renew_flag"); v != nil {
		request.RenewFlag = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().RenewDataEngine(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate dlc renewDataEngine failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dataEngineName)
	return resourceTencentCloudDlcRenewDataEngineRead(d, meta)
}

func resourceTencentCloudDlcRenewDataEngineRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_renew_data_engine_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcRenewDataEngineDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_renew_data_engine_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

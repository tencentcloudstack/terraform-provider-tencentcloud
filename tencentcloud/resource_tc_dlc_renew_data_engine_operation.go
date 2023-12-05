package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDlcRenewDataEngineOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcRenewDataEngineCreate,
		Read:   resourceTencentCloudDlcRenewDataEngineRead,
		Delete: resourceTencentCloudDlcRenewDataEngineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"data_engine_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Data engine name.",
			},

			"time_span": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Engine TimeSpan, prePay: minimum of 1, representing one month of purchasing resources, with a maximum of 120, default 3600, postPay: fixed fee of 3600.",
			},

			"pay_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Engine pay mode type, only support 0: postPay, 1: prePay(default).",
			},

			"time_unit": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Engine TimeUnit, prePay: use m(default), postPay: use h.",
			},

			"renew_flag": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Automatic renewal flag, 0, initial state, automatic renewal is not performed by default. if the user has prepaid non-stop service privileges, automatic renewal will occur. 1: Automatic renewal. 2: make it clear that there will be no automatic renewal. if this parameter is not passed, the default value is 0.",
			},
		},
	}
}

func resourceTencentCloudDlcRenewDataEngineCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_renew_data_engine_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().RenewDataEngine(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_dlc_renew_data_engine_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcRenewDataEngineDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_renew_data_engine_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

package cwp

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cwpv20180228 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCwpAutoOpenProversionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCwpAutoOpenProversionConfigCreate,
		Read:   resourceTencentCloudCwpAutoOpenProversionConfigRead,
		Update: resourceTencentCloudCwpAutoOpenProversionConfigUpdate,
		Delete: resourceTencentCloudCwpAutoOpenProversionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Set the auto-activation status.\n<li>CLOSE: off</li>\n<li>OPEN: on</li>.",
			},

			"protect_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enhanced Protection Mode PROVERSION_POSTPAY Professional Edition - Pay-as-you-go PROVERSION_PREPAY Professional Edition - Annual/Monthly Subscription FLAGSHIP_PREPAY Flagship Edition - Annual/Monthly Subscription.",
			},

			"auto_repurchase_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Automatic purchase/expansion authorization switch, 1 by default, 0 for OFF, 1 for ON.",
			},

			"auto_repurchase_renew_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Auto-renewal or not for auto-purchased orders, 0 by default, 0 for OFF, 1 for ON.",
			},

			"repurchase_renew_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether the manually purchased order is automatically renewed (defaults to 0). 0 - off; 1 -on.",
			},

			"auto_bind_rasp_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Newly added machines will be automatically bound to Rasp. 0: Disabled, 1: Enabled.",
			},

			"auto_open_rasp_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Newly added machines will have automatic Raspberry Pi protection enabled by default. (0: Disabled, 1: Enabled).",
			},

			"auto_downgrade_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Automatic scaling switch: 0 for off, 1 for on.",
			},
		},
	}
}

func resourceTencentCloudCwpAutoOpenProversionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cwp_auto_open_proversion_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudCwpAutoOpenProversionConfigUpdate(d, meta)
}

func resourceTencentCloudCwpAutoOpenProversionConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cwp_auto_open_proversion_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CwpService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeCwpAutoOpenProversionConfigById(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cwp_auto_open_proversion_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.AutoOpenStatus != nil {
		if *respData.AutoOpenStatus == true {
			_ = d.Set("status", "OPEN")
		} else {
			_ = d.Set("status", "CLOSE")
		}
	}

	if respData.ProtectType != nil {
		_ = d.Set("protect_type", respData.ProtectType)
	}

	if respData.AutoRepurchaseSwitch != nil {
		if *respData.AutoRepurchaseSwitch == true {
			_ = d.Set("auto_repurchase_switch", 1)
		} else {
			_ = d.Set("auto_repurchase_switch", 0)
		}
	}

	if respData.AutoRepurchaseRenewSwitch != nil {
		if *respData.AutoRepurchaseRenewSwitch == true {
			_ = d.Set("auto_repurchase_renew_switch", 1)
		} else {
			_ = d.Set("auto_repurchase_renew_switch", 0)
		}
	}

	if respData.RepurchaseRenewSwitch != nil {
		if *respData.RepurchaseRenewSwitch == true {
			_ = d.Set("repurchase_renew_switch", 1)
		} else {
			_ = d.Set("repurchase_renew_switch", 0)
		}
	}

	if respData.AutoBindRaspSwitch != nil {
		if *respData.AutoBindRaspSwitch == true {
			_ = d.Set("auto_bind_rasp_switch", 1)
		} else {
			_ = d.Set("auto_bind_rasp_switch", 0)
		}
	}

	if respData.AutoOpenRaspSwitch != nil {
		if *respData.AutoOpenRaspSwitch == true {
			_ = d.Set("auto_open_rasp_switch", 1)
		} else {
			_ = d.Set("auto_open_rasp_switch", 0)
		}
	}

	if respData.AutoDowngradeSwitch != nil {
		if *respData.AutoDowngradeSwitch == true {
			_ = d.Set("auto_downgrade_switch", 1)
		} else {
			_ = d.Set("auto_downgrade_switch", 0)
		}
	}

	return nil
}

func resourceTencentCloudCwpAutoOpenProversionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cwp_auto_open_proversion_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cwpv20180228.NewModifyAutoOpenProVersionConfigRequest()
	)

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protect_type"); ok {
		request.ProtectType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_repurchase_switch"); ok {
		request.AutoRepurchaseSwitch = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_repurchase_renew_switch"); ok {
		request.AutoRepurchaseRenewSwitch = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("repurchase_renew_switch"); ok {
		request.RepurchaseRenewSwitch = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_bind_rasp_switch"); ok {
		request.AutoBindRaspSwitch = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_open_rasp_switch"); ok {
		request.AutoOpenRaspSwitch = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_downgrade_switch"); ok {
		request.AutoDowngradeSwitch = helper.IntUint64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCwpClient().ModifyAutoOpenProVersionConfigWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update cwp auto open proversion config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudCwpAutoOpenProversionConfigRead(d, meta)
}

func resourceTencentCloudCwpAutoOpenProversionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cwp_auto_open_proversion_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

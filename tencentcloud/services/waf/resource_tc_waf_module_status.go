package waf

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafModuleStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafModuleStatusCreate,
		Read:   resourceTencentCloudWafModuleStatusRead,
		Update: resourceTencentCloudWafModuleStatusUpdate,
		Delete: resourceTencentCloudWafModuleStatusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},
			"web_security": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "WEB security module status, 0:closed, 1:opened.",
			},
			"access_control": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "ACL module status, 0:closed, 1:opened.",
			},
			"cc_protection": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "CC module status, 0:closed, 1:opened.",
			},
			"api_protection": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "API security module status, 0:closed, 1:opened.",
			},
			"anti_tamper": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "Anti tamper module status, 0:closed, 1:opened.",
			},
			"anti_leakage": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "Anti leakage module status, 0:closed, 1:opened.",
			},
		},
	}
}

func resourceTencentCloudWafModuleStatusCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_module_status.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var domain string
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	d.SetId(domain)

	return resourceTencentCloudWafModuleStatusUpdate(d, meta)
}

func resourceTencentCloudWafModuleStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_module_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		domain  = d.Id()
	)

	moduleStatus, err := service.DescribeWafModuleStatusById(ctx, domain)
	if err != nil {
		return err
	}

	if moduleStatus == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafModuleStatus` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)

	if moduleStatus.WebSecurity != nil {
		_ = d.Set("web_security", moduleStatus.WebSecurity)
	}

	if moduleStatus.AccessControl != nil {
		_ = d.Set("access_control", moduleStatus.AccessControl)
	}

	if moduleStatus.CcProtection != nil {
		_ = d.Set("cc_protection", moduleStatus.CcProtection)
	}

	if moduleStatus.ApiProtection != nil {
		_ = d.Set("api_protection", moduleStatus.ApiProtection)
	}

	if moduleStatus.AntiTamper != nil {
		_ = d.Set("anti_tamper", moduleStatus.AntiTamper)
	}

	if moduleStatus.AntiLeakage != nil {
		_ = d.Set("anti_leakage", moduleStatus.AntiLeakage)
	}

	return nil
}

func resourceTencentCloudWafModuleStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_module_status.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = waf.NewModifyModuleStatusRequest()
		domain  = d.Id()
	)

	request.Domain = &domain

	if v, ok := d.GetOkExists("web_security"); ok {
		request.WebSecurity = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("access_control"); ok {
		request.AccessControl = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("cc_protection"); ok {
		request.CcProtection = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("api_protection"); ok {
		request.ApiProtection = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("anti_tamper"); ok {
		request.AntiTamper = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("anti_leakage"); ok {
		request.AntiLeakage = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyModuleStatus(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("waf moduleStatus version not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify waf moduleStatus failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWafModuleStatusRead(d, meta)
}

func resourceTencentCloudWafModuleStatusDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_module_status.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

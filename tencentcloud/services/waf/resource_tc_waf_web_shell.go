package waf

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafWebShell() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafWebShellCreate,
		Read:   resourceTencentCloudWafWebShellRead,
		Update: resourceTencentCloudWafWebShellUpdate,
		Delete: resourceTencentCloudWafWebShellDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain.",
			},
			"status": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1, 2}),
				Description:  "Webshell status, 1: open; 0: closed; 2: log.",
			},
		},
	}
}

func resourceTencentCloudWafWebShellCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_web_shell.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var domain string

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	d.SetId(domain)

	return resourceTencentCloudWafWebShellUpdate(d, meta)
}

func resourceTencentCloudWafWebShellRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_web_shell.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		domain  = d.Id()
	)

	webShell, err := service.DescribeWafWebShellById(ctx, domain)
	if err != nil {
		return err
	}

	if webShell == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafWebShell` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if webShell.Domain != nil {
		_ = d.Set("domain", webShell.Domain)
	}

	if webShell.Status != nil {
		_ = d.Set("status", webShell.Status)
	}

	return nil
}

func resourceTencentCloudWafWebShellUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_web_shell.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = waf.NewModifyWebshellStatusRequest()
		domain  = d.Id()
	)

	webShellStatus := waf.WebshellStatus{}
	webShellStatus.Domain = helper.String(domain)

	if v, ok := d.GetOkExists("status"); ok {
		webShellStatus.Status = helper.IntUint64(v.(int))
	}

	request.Webshell = &webShellStatus

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyWebshellStatus(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update waf webShell failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWafWebShellRead(d, meta)
}

func resourceTencentCloudWafWebShellDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_web_shell.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

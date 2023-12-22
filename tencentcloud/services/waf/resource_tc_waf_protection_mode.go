package waf

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafProtectionMode() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafProtectionModeCreate,
		Read:   resourceTencentCloudWafProtectionModeRead,
		Update: resourceTencentCloudWafProtectionModeUpdate,
		Delete: resourceTencentCloudWafProtectionModeDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},
			"mode": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Protection status:10: Rule observation; AI off mode, 11: Rule observation; AI observation mode, 12: Rule observation; AI interception mode20: Rule interception; AI off mode, 21: Rule interception; AI observation mode, 22: Rule interception; AI interception mode.",
			},
			"edition": {
				Optional:     true,
				Type:         schema.TypeString,
				Default:      "sparta-waf",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"clb-waf", "sparta-waf"}),
				Description:  "WAF edition. clb-waf means clb-waf, sparta-waf means saas-waf, default is sparta-waf.",
			},
			"type": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "0 is to modify the rule engine status, 1 is to modify the AI status.",
			},
		},
	}
}

func resourceTencentCloudWafProtectionModeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_protection_mode.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		domain  string
		edition string
	)

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	if v, ok := d.GetOk("edition"); ok {
		edition = v.(string)
	}

	d.SetId(strings.Join([]string{domain, edition}, tccommon.FILED_SP))

	return resourceTencentCloudWafProtectionModeUpdate(d, meta)
}

func resourceTencentCloudWafProtectionModeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_protection_mode.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domain := idSplit[0]
	edition := idSplit[1]

	protectionInfo, err := service.DescribeSpartaProtectionInfoById(ctx, domain, edition)
	if err != nil {
		return err
	}

	if protectionInfo == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafProtectionMode` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("edition", edition)

	if protectionInfo.Engine != nil {
		engineIne, _ := strconv.Atoi(*protectionInfo.Engine)
		_ = d.Set("mode", engineIne)
	}

	return nil
}

func resourceTencentCloudWafProtectionModeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_protection_mode.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = waf.NewModifySpartaProtectionModeRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domain := idSplit[0]
	edition := idSplit[1]

	request.Domain = &domain
	request.Edition = &edition

	if v, _ := d.GetOkExists("mode"); v != nil {
		request.Mode = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOkExists("type"); v != nil {
		request.Type = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifySpartaProtectionMode(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate waf protectionMode failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWafProtectionModeRead(d, meta)
}

func resourceTencentCloudWafProtectionModeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_protection_mode.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

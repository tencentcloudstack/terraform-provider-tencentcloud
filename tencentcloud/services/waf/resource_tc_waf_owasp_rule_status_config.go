package waf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafOwaspRuleStatusConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafOwaspRuleStatusConfigCreate,
		Read:   resourceTencentCloudWafOwaspRuleStatusConfigRead,
		Update: resourceTencentCloudWafOwaspRuleStatusConfigUpdate,
		Delete: resourceTencentCloudWafOwaspRuleStatusConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name.",
			},

			"rule_status": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Rule switch. valid values: 0 (disabled), 1 (enabled), 2 (observation only).",
			},

			"rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule ID.",
			},

			"type_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "If reverse requires the input of data type.",
			},

			"reason": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Reason for modification. valid values: 0: none (compatibility record is empty). 1: avoid false positives due to business characteristics. 2: reporting of rule-based false positives. 3: gray release of core business rules. 4: others.",
			},

			// computed
			"cve_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CVE ID.",
			},

			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule description.",
			},

			"level": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Protection level of the rule. valid values: 100 (loose), 200 (normal), 300 (strict), 400 (ultra-strict).",
			},

			"vul_level": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Threat level. valid values: 0 (unknown), 100 (low risk), 200 (medium risk), 300 (high risk), 400 (critical).",
			},

			"locked": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the user is locked.",
			},
		},
	}
}

func resourceTencentCloudWafOwaspRuleStatusConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_rule_status_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		domain string
		ruleId string
	)

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	if v, ok := d.GetOk("rule_id"); ok {
		ruleId = v.(string)
	}

	d.SetId(strings.Join([]string{domain, ruleId}, tccommon.FILED_SP))
	return resourceTencentCloudWafOwaspRuleStatusConfigUpdate(d, meta)
}

func resourceTencentCloudWafOwaspRuleStatusConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_rule_status_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	respData, err := service.DescribeWafOwaspRuleStatusConfigById(ctx, domain, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_waf_owasp_rule_status_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("rule_id", ruleId)

	if respData.Status != nil {
		_ = d.Set("rule_status", respData.Status)
	}

	if respData.TypeId != nil {
		_ = d.Set("type_id", respData.TypeId)
	}

	if respData.Reason != nil {
		_ = d.Set("reason", respData.Reason)
	}

	if respData.CveID != nil {
		_ = d.Set("cve_id", respData.CveID)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Level != nil {
		_ = d.Set("level", respData.Level)
	}

	if respData.VulLevel != nil {
		_ = d.Set("vul_level", respData.VulLevel)
	}

	if respData.Locked != nil {
		_ = d.Set("locked", respData.Locked)
	}

	return nil
}

func resourceTencentCloudWafOwaspRuleStatusConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_rule_status_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	request := wafv20180125.NewModifyOwaspRuleStatusRequest()
	if v, ok := d.GetOkExists("rule_status"); ok {
		request.RuleStatus = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("type_id"); ok {
		request.TypeId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("reason"); ok {
		request.Reason = helper.IntInt64(v.(int))
	}

	request.Domain = &domain
	request.RuleIDs = append(request.RuleIDs, &ruleId)
	request.SelectAll = helper.Bool(false)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyOwaspRuleStatusWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update waf owasp rule status config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudWafOwaspRuleStatusConfigRead(d, meta)
}

func resourceTencentCloudWafOwaspRuleStatusConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_rule_status_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

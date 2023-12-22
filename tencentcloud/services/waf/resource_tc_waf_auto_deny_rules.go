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

func ResourceTencentCloudWafAutoDenyRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafAutoDenyRulesCreate,
		Read:   resourceTencentCloudWafAutoDenyRulesRead,
		Delete: resourceTencentCloudWafAutoDenyRulesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},
			"attack_threshold": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 100),
				Description:  "The threshold number of attacks that triggers IP autodeny, ranging from 2 to 100 times.",
			},
			"time_threshold": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 60),
				Description:  "IP autodeny statistical time, ranging from 1-60 minutes.",
			},
			"deny_time_threshold": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerInRange(5, 360),
				Description:  "The IP autodeny time after triggering the IP autodeny, ranging from 5 to 360 minutes.",
			},
		},
	}
}

func resourceTencentCloudWafAutoDenyRulesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_auto_deny_rules.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = waf.NewModifyWafAutoDenyRulesRequest()
		domain  string
	)

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOkExists("attack_threshold"); ok {
		request.AttackThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("time_threshold"); ok {
		request.TimeThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("deny_time_threshold"); ok {
		request.DenyTimeThreshold = helper.IntInt64(v.(int))
	}

	request.DefenseStatus = helper.IntInt64(1)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyWafAutoDenyRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || *result.Response.Success.Code != "Success" {
			e = fmt.Errorf("create waf autoDenyRules not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create waf autoDenyRules failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(domain)

	return resourceTencentCloudWafAutoDenyRulesRead(d, meta)
}

func resourceTencentCloudWafAutoDenyRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_auto_deny_rules.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		domain  = d.Id()
	)

	autoDenyRules, err := service.DescribeWafAutoDenyRulesById(ctx, domain)
	if err != nil {
		return err
	}

	if autoDenyRules == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafAutoDenyRules` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)

	if autoDenyRules.AttackThreshold != nil {
		_ = d.Set("attack_threshold", autoDenyRules.AttackThreshold)
	}

	if autoDenyRules.TimeThreshold != nil {
		_ = d.Set("time_threshold", autoDenyRules.TimeThreshold)
	}

	if autoDenyRules.DenyTimeThreshold != nil {
		_ = d.Set("deny_time_threshold", autoDenyRules.DenyTimeThreshold)
	}

	return nil
}

func resourceTencentCloudWafAutoDenyRulesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_auto_deny_rules.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = waf.NewModifyWafAutoDenyRulesRequest()
		domain  = d.Id()
	)

	if v, ok := d.GetOkExists("attack_threshold"); ok {
		request.AttackThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("time_threshold"); ok {
		request.TimeThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("deny_time_threshold"); ok {
		request.DenyTimeThreshold = helper.IntInt64(v.(int))
	}

	request.Domain = &domain
	request.DefenseStatus = helper.IntInt64(0)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyWafAutoDenyRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || *result.Response.Success.Code != "Success" {
			e = fmt.Errorf("delete waf autoDenyRules not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete waf autoDenyRules failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

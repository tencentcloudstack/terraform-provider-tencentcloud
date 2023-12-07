package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWafProtectionMode() *schema.Resource {
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
				ValidateFunc: validateAllowedStringValue([]string{"clb-waf", "sparta-waf"}),
				Description:  "WAF edition. clb-waf means clb-waf, sparta-waf means saas-waf, default is sparta-waf.",
			},
			"type": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "0 is to modify the rule engine status, 1 is to modify the AI status.",
			},
		},
	}
}

func resourceTencentCloudWafProtectionModeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_protection_mode.create")()
	defer inconsistentCheck(d, meta)()

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

	d.SetId(strings.Join([]string{domain, edition}, FILED_SP))

	return resourceTencentCloudWafProtectionModeUpdate(d, meta)
}

func resourceTencentCloudWafProtectionModeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_protection_mode.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
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
	defer logElapsed("resource.tencentcloud_waf_protection_mode.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = waf.NewModifySpartaProtectionModeRequest()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifySpartaProtectionMode(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_waf_protection_mode.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

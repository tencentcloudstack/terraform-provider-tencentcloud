package waf

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafDomainPostActionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafDomainPostActionConfigCreate,
		Read:   resourceTencentCloudWafDomainPostActionConfigRead,
		Update: resourceTencentCloudWafDomainPostActionConfigUpdate,
		Delete: resourceTencentCloudWafDomainPostActionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain.",
			},

			"post_cls_action": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "0- Disable shipping, 1- Enable shipping.",
			},

			"post_ckafka_action": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "0- Disable shipping, 1- Enable shipping.",
			},
		},
	}
}

func resourceTencentCloudWafDomainPostActionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_domain_post_action_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var domain string
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	d.SetId(domain)

	return resourceTencentCloudWafDomainPostActionConfigUpdate(d, meta)
}

func resourceTencentCloudWafDomainPostActionConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_domain_post_action_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		domain  = d.Id()
	)

	respData, err := service.DescribeWafDomainPostActionById(ctx, domain)
	if err != nil {
		return err
	}

	if respData == nil || len(respData) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `waf_domain_post_action_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)

	for _, item := range respData {
		if item.PostCLSStatus != nil {
			_ = d.Set("post_cls_action", item.PostCLSStatus)
		}

		if item.PostCKafkaStatus != nil {
			_ = d.Set("post_ckafka_action", item.PostCKafkaStatus)
		}
	}

	return nil
}

func resourceTencentCloudWafDomainPostActionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_domain_post_action_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wafv20180125.NewModifyDomainPostActionRequest()
		domain  = d.Id()
	)

	if v, ok := d.GetOkExists("post_cls_action"); ok {
		request.PostCLSAction = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("post_ckafka_action"); ok {
		request.PostCKafkaAction = helper.IntInt64(v.(int))
	}

	request.Domain = &domain
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyDomainPostActionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update waf domain post action config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudWafDomainPostActionConfigRead(d, meta)
}

func resourceTencentCloudWafDomainPostActionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_domain_post_action_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

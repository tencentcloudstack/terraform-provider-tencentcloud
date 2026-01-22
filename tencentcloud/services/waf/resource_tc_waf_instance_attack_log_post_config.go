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

func ResourceTencentCloudWafInstanceAttackLogPostConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafInstanceAttackLogPostConfigCreate,
		Read:   resourceTencentCloudWafInstanceAttackLogPostConfigRead,
		Update: resourceTencentCloudWafInstanceAttackLogPostConfigUpdate,
		Delete: resourceTencentCloudWafInstanceAttackLogPostConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Waf instance ID.",
			},

			"attack_log_post": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "Attack log delivery switch. 0- Disable, 1- Enable.",
			},
		},
	}
}

func resourceTencentCloudWafInstanceAttackLogPostConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_instance_attack_log_post_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudWafInstanceAttackLogPostConfigUpdate(d, meta)
}

func resourceTencentCloudWafInstanceAttackLogPostConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_instance_attack_log_post_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	respData, err := service.DescribeWafInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `waf_instance_attack_log_post_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if respData.AttackLogPost != nil {
		_ = d.Set("attack_log_post", respData.AttackLogPost)
	}

	return nil
}

func resourceTencentCloudWafInstanceAttackLogPostConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_instance_attack_log_post_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = wafv20180125.NewModifyInstanceAttackLogPostRequest()
		instanceId = d.Id()
	)

	if v, ok := d.GetOkExists("attack_log_post"); ok {
		request.AttackLogPost = helper.IntInt64(v.(int))
	}

	request.InstanceId = &instanceId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyInstanceAttackLogPostWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update waf instance attack log post config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudWafInstanceAttackLogPostConfigRead(d, meta)
}

func resourceTencentCloudWafInstanceAttackLogPostConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_instance_attack_log_post_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

package waf

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafModifyAccessPeriod() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafModifyAccessPeriodCreate,
		Read:   resourceTencentCloudWafModifyAccessPeriodRead,
		Delete: resourceTencentCloudWafModifyAccessPeriodDelete,

		Schema: map[string]*schema.Schema{
			"period": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 180),
				Description:  "Access log retention period, range is [1, 180].",
			},
			"topic_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Log topic, new version does not need to be uploaded.",
			},
		},
	}
}

func resourceTencentCloudWafModifyAccessPeriodCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_modify_access_period.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = waf.NewModifyAccessPeriodRequest()
		topicId string
	)

	if v, _ := d.GetOkExists("period"); v != nil {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
		topicId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyAccessPeriod(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate waf ModifyAccessPeriod failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(topicId)

	return resourceTencentCloudWafModifyAccessPeriodRead(d, meta)
}

func resourceTencentCloudWafModifyAccessPeriodRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_modify_access_period.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudWafModifyAccessPeriodDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_modify_access_period.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

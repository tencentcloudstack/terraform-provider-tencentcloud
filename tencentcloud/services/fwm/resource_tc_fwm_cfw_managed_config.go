package fwm

import (
	"context"
	"log"

	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudFwmCfwManagedConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmCfwManagedConfigCreate,
		Read:   resourceTencentCloudFwmCfwManagedConfigRead,
		Update: resourceTencentCloudFwmCfwManagedConfigUpdate,
		Delete: resourceTencentCloudFwmCfwManagedConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"member_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target member id of the organization.",
			},
			"cfw_managed": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Cloud firewall managed switch, 1: enable managed, 0: disable managed.",
			},
		},
	}
}

func resourceTencentCloudFwmCfwManagedConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_managed_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	memberId := d.Get("member_id").(string)
	d.SetId(memberId)

	return resourceTencentCloudFwmCfwManagedConfigUpdate(d, meta)
}

func resourceTencentCloudFwmCfwManagedConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_managed_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		memberId = d.Id()
	)

	member, err := service.DescribeFwmMemberGroup(ctx, memberId)
	if err != nil {
		return err
	}

	if member == nil {
		log.Printf("[WARN]%s resource `tencentcloud_fwm_cfw_managed_config` id=%s not found, maybe the member has left the organization.", logId, d.Id())
		d.SetId("")
		return nil
	}

	if member.MemberId != nil {
		_ = d.Set("member_id", member.MemberId)
	}

	if member.CfwManaged != nil {
		_ = d.Set("cfw_managed", member.CfwManaged)
	}

	return nil
}

func resourceTencentCloudFwmCfwManagedConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_managed_config.update")()
	defer tccommon.InconsistentCheck(d, meta)

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		memberId = d.Id()
	)

	request := fwmv20250611.NewModifyCfwManagedRequest()
	if v, ok := d.GetOkExists("cfw_managed"); ok {
		request.CfwManaged = helper.IntInt64(v.(int))
	}

	request.MemberId = helper.String(memberId)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmClient().ModifyCfwManaged(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update fwm cfw managed failed, reason:%+v", logId, err)
	}

	return resourceTencentCloudFwmCfwManagedConfigRead(d, meta)
}

func resourceTencentCloudFwmCfwManagedConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_managed_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)

	return nil
}

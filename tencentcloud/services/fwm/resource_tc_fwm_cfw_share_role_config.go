package fwm

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudFwmCfwShareRoleConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmCfwShareRoleConfigConfigCreate,
		Read:   resourceTencentCloudFwmCfwShareRoleConfigConfigRead,
		Update: resourceTencentCloudFwmCfwShareRoleConfigConfigUpdate,
		Delete: resourceTencentCloudFwmCfwShareRoleConfigConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cfw_share_role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud firewall shared roles, with values of sharer, user, and none. Sharer -> user prohibited from direct changes, must first be set to none.",
			},
			"user_member_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User ID.",
			},
			"sharer_member_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sharer account ID, required when CfwShareRole is user.",
			},
		},
	}
}

func resourceTencentCloudFwmCfwShareRoleConfigConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_share_role_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	userMemberId := d.Get("user_member_id").(string)

	d.SetId(userMemberId)

	return resourceTencentCloudFwmCfwShareRoleConfigConfigUpdate(d, meta)
}

func resourceTencentCloudFwmCfwShareRoleConfigConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_share_role_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service      = FwmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		userMemberId = d.Id()
	)

	organMember, err := service.DescribeFwmMemberGroup(ctx, userMemberId)
	if err != nil {
		return err
	}

	if organMember == nil {
		log.Printf("[WARN]%s resource `tencentcloud_fwm_cfw_share_role_config` [%s] member not found", logId, userMemberId)
		d.SetId("")
		return nil
	}

	if organMember.CfwShareRole != nil {
		_ = d.Set("cfw_share_role", organMember.CfwShareRole)
	}

	if organMember.MemberId != nil {
		_ = d.Set("user_member_id", organMember.MemberId)
	}

	return nil
}

func resourceTencentCloudFwmCfwShareRoleConfigConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_share_role_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request      = fwm.NewModifyCfwShareRoleRequest()
		userMemberId = d.Id()
	)

	if v, ok := d.GetOk("cfw_share_role"); ok {
		request.CfwShareRole = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sharer_member_id"); ok {
		request.SharerMemberId = helper.String(v.(string))
	}

	request.UserMemberId = helper.String(userMemberId)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmClient().ModifyCfwShareRoleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify fwm cfw share role to %s failed, reason:%+v", logId, userMemberId, err)
		return err
	}

	return resourceTencentCloudFwmCfwShareRoleConfigConfigRead(d, meta)
}

func resourceTencentCloudFwmCfwShareRoleConfigConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_share_role_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

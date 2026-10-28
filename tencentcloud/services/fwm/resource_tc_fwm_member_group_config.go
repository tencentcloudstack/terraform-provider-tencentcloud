package fwm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"
)

func ResourceTencentCloudFwmMemberGroupConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmMemberGroupConfigCreate,
		Read:   resourceTencentCloudFwmMemberGroupConfigRead,
		Update: resourceTencentCloudFwmMemberGroupConfigUpdate,
		Delete: resourceTencentCloudFwmMemberGroupConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_member_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Member ID of the FWM organization member, such as `mem-xxxx`.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account group ID that the member belongs to. Empty string means the member is not grouped (moved out). A member can belong to at most one group.",
			},
		},
	}
}

func resourceTencentCloudFwmMemberGroupConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_member_group_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	userMemberId := d.Get("user_member_id").(string)
	d.SetId(userMemberId)

	return resourceTencentCloudFwmMemberGroupConfigUpdate(d, meta)
}

func resourceTencentCloudFwmMemberGroupConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_member_group_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service      = FwmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		userMemberId = d.Id()
	)

	var memberItem *fwmv20250611.OrganMemberItem
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		item, e := service.DescribeFwmMemberGroup(ctx, userMemberId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		memberItem = item
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read fwm_member_group failed, id=%s, reason:%s", logId, userMemberId, err.Error())
		return err
	}

	if memberItem == nil {
		log.Printf("[WARN]%s resource `tencentcloud_fwm_member_group_config` [%s] member not found", logId, userMemberId)
		d.SetId("")
		return nil
	}

	if memberItem.MemberId != nil {
		_ = d.Set("user_member_id", *memberItem.MemberId)
	}

	groupId := ""
	if memberItem.AccountGroup != nil && memberItem.AccountGroup.GroupId != nil {
		groupId = *memberItem.AccountGroup.GroupId
	}
	_ = d.Set("group_id", groupId)

	return nil
}

func resourceTencentCloudFwmMemberGroupConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_member_group_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service      = FwmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		userMemberId = d.Id()
	)

	if d.HasChange("group_id") {
		groupId := d.Get("group_id").(string)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := service.ModifyFwmMemberGroup(ctx, userMemberId, groupId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update fwm_member_group failed, id=%s, reason:%s", logId, userMemberId, err.Error())
			return err
		}
	}

	return resourceTencentCloudFwmMemberGroupConfigRead(d, meta)
}

func resourceTencentCloudFwmMemberGroupConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_member_group_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

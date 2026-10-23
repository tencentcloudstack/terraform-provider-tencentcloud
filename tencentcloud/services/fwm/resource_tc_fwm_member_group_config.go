package fwm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"
)

func ResourceTencentCloudFwmMemberGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmMemberGroupCreate,
		Read:   resourceTencentCloudFwmMemberGroupRead,
		Update: resourceTencentCloudFwmMemberGroupUpdate,
		Delete: resourceTencentCloudFwmMemberGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"member_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Member ID of the FWM organization member, such as `mem-xxxx`.",
			},
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "AppId of the FWM organization member (numeric string), used to locate the member because the cloud API filter does not support MemberId.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account group ID that the member belongs to. Empty string means the member is not grouped (moved out). A member can belong to at most one group.",
			},
		},
	}
}

func resourceTencentCloudFwmMemberGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_member_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	appId := d.Get("app_id").(string)
	memberId := d.Get("member_id").(string)

	id := appId + tccommon.FILED_SP + memberId
	d.SetId(id)

	return resourceTencentCloudFwmMemberGroupUpdate(d, meta)
}

func resourceTencentCloudFwmMemberGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_member_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = FwmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id      = d.Id()
	)

	appId, memberId, err := service.FwmMemberGroupIdToAppIdMemberId(id)
	if err != nil {
		log.Printf("[CRITAL]%s parse fwm_member_group id failed, id=%s, reason:%s", logId, id, err.Error())
		return err
	}

	var memberItem *fwm.OrganMemberItem
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		item, e := service.DescribeFwmMemberGroup(ctx, appId, memberId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		memberItem = item
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read fwm_member_group failed, id=%s, reason:%s", logId, id, err.Error())
		return err
	}

	if memberItem == nil {
		log.Printf("[WARN]%s resource `fwm_member_group` [%s] member not found", logId, id)
		d.SetId("")
		return nil
	}

	groupId := ""
	if memberItem.AccountGroup != nil && memberItem.AccountGroup.GroupId != nil {
		groupId = *memberItem.AccountGroup.GroupId
	}
	_ = d.Set("group_id", groupId)

	if memberItem.AppId != nil {
		_ = d.Set("app_id", *memberItem.AppId)
	}
	if memberItem.MemberId != nil {
		_ = d.Set("member_id", *memberItem.MemberId)
	}

	return nil
}

func resourceTencentCloudFwmMemberGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_member_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = FwmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id      = d.Id()
	)

	_, memberId, err := service.FwmMemberGroupIdToAppIdMemberId(id)
	if err != nil {
		log.Printf("[CRITAL]%s parse fwm_member_group id failed, id=%s, reason:%s", logId, id, err.Error())
		return err
	}

	if d.HasChange("group_id") {
		groupId := d.Get("group_id").(string)

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := service.ModifyFwmMemberGroup(ctx, memberId, groupId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update fwm_member_group failed, id=%s, reason:%s", logId, id, err.Error())
			return err
		}
	}

	return resourceTencentCloudFwmMemberGroupRead(d, meta)
}

func resourceTencentCloudFwmMemberGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_member_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

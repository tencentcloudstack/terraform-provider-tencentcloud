package tco

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRejectJoinShareUnitInvitationOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRejectJoinShareUnitInvitationOperationCreate,
		Read:   resourceTencentCloudRejectJoinShareUnitInvitationOperationRead,
		Delete: resourceTencentCloudRejectJoinShareUnitInvitationOperationDelete,
		Schema: map[string]*schema.Schema{
			"unit_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Shared unit ID.",
			},
		},
	}
}

func resourceTencentCloudRejectJoinShareUnitInvitationOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_reject_join_share_unit_invitation_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		unitId string
	)
	var (
		request  = organization.NewRejectJoinShareUnitInvitationRequest()
		response = organization.NewRejectJoinShareUnitInvitationResponse()
	)

	if v, ok := d.GetOk("unit_id"); ok {
		request.UnitId = helper.String(v.(string))
		unitId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().RejectJoinShareUnitInvitationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create reject join share unit invitation operation failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	d.SetId(unitId)

	return resourceTencentCloudRejectJoinShareUnitInvitationOperationRead(d, meta)
}

func resourceTencentCloudRejectJoinShareUnitInvitationOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_reject_join_share_unit_invitation_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRejectJoinShareUnitInvitationOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_reject_join_share_unit_invitation_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

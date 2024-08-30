package tco

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIdentityCenterUserGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIdentityCenterUserGroupAttachmentCreate,
		Read:   resourceTencentCloudIdentityCenterUserGroupAttachmentRead,
		Delete: resourceTencentCloudIdentityCenterUserGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone id.",
			},

			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User group ID.",
			},

			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User ID.",
			},
		},
	}
}

func resourceTencentCloudIdentityCenterUserGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_user_group_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId  string
		groupId string
		userId  string
	)
	var (
		request  = organization.NewAddUserToGroupRequest()
		response = organization.NewAddUserToGroupResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
	}
	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_id"); ok {
		request.UserId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().AddUserToGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create identity center user group attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	d.SetId(strings.Join([]string{zoneId, groupId, userId}, tccommon.FILED_SP))

	return resourceTencentCloudIdentityCenterUserGroupAttachmentRead(d, meta)
}

func resourceTencentCloudIdentityCenterUserGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_user_group_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	groupId := idSplit[1]
	userId := idSplit[2]

	respData, err := service.DescribeIdentityCenterUserGroupAttachmentById(ctx, zoneId, groupId, userId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `identity_center_user_group_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.GroupId != nil {
		_ = d.Set("group_id", respData.GroupId)
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("user_id", userId)
	return nil
}

func resourceTencentCloudIdentityCenterUserGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_user_group_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	groupId := idSplit[1]
	userId := idSplit[2]

	var (
		request  = organization.NewRemoveUserFromGroupRequest()
		response = organization.NewRemoveUserFromGroupResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_id"); ok {
		request.UserId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().RemoveUserFromGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete identity center user group attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = zoneId
	_ = groupId
	_ = userId
	return nil
}

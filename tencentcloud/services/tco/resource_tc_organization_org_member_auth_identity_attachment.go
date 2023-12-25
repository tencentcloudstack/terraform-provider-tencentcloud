package tco

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOrganizationOrgMemberAuthIdentityAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgMemberAuthIdentityAttachmentCreate,
		Read:   resourceTencentCloudOrganizationOrgMemberAuthIdentityAttachmentRead,
		Delete: resourceTencentCloudOrganizationOrgMemberAuthIdentityAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"member_uin": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Member Uin.",
			},

			"identity_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Identity Id list. Up to 5.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgMemberAuthIdentityAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_member_auth_identity.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = organization.NewCreateOrganizationMemberAuthIdentityRequest()
		memberUin string
	)
	if v, ok := d.GetOk("member_uin"); ok {
		request.MemberUins = append(request.MemberUins, helper.IntUint64(v.(int)))
		memberUin = helper.IntToStr(v.(int))
	}

	if v, ok := d.GetOk("identity_ids"); ok {
		identityIdsSet := v.(*schema.Set).List()
		for i := range identityIdsSet {
			identityIds := identityIdsSet[i].(int)
			request.IdentityIds = append(request.IdentityIds, helper.IntUint64(identityIds))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().CreateOrganizationMemberAuthIdentity(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization orgMemberAuthIdentity failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(memberUin)
	return resourceTencentCloudOrganizationOrgMemberAuthIdentityAttachmentRead(d, meta)
}

func resourceTencentCloudOrganizationOrgMemberAuthIdentityAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_member_auth_identity.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	memberUin := d.Id()
	uin := helper.StrToInt64(memberUin)
	identityIds, err := service.DescribeOrganizationOrgMemberAuthIdentityById(ctx, uin)
	if err != nil {
		return err
	}

	if len(identityIds) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgMemberAuthIdentity` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("identity_ids", identityIds)
	_ = d.Set("member_uin", uin)
	return nil
}

func resourceTencentCloudOrganizationOrgMemberAuthIdentityAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_member_auth_identity.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	uin := d.Id()
	var identityIds []string
	if v, ok := d.GetOk("identity_ids"); ok {
		identityIdsSet := v.(*schema.Set).List()
		for i := range identityIdsSet {
			identityId := identityIdsSet[i].(int)
			identityIds = append(identityIds, helper.IntToStr(identityId))
		}
	}

	if err := service.DeleteOrganizationOrgMemberAuthIdentityById(ctx, uin, identityIds); err != nil {
		return err
	}

	return nil
}

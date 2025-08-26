package tco

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organizationv20210331 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOrganizationMemberAuthPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationMemberAuthPolicyAttachmentCreate,
		Read:   resourceTencentCloudOrganizationMemberAuthPolicyAttachmentRead,
		Delete: resourceTencentCloudOrganizationMemberAuthPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Pilicy ID.",
			},

			"org_sub_account_uin": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Organization administrator sub-account Uin.",
			},

			// computed
			"identity_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Identity ID.",
			},

			"identity_role_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity role name.",
			},

			"identity_role_alias_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity role alias name.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time.",
			},

			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy name.",
			},

			"member_uin": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Member UIN.",
			},

			"member_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Member name.",
			},

			"org_sub_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Org sub account name.",
			},

			"bind_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Bind type. 1-Subaccount, 2-User Group.",
			},
		},
	}
}

func resourceTencentCloudOrganizationMemberAuthPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_member_auth_policy_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request          = organizationv20210331.NewBindOrganizationPolicySubAccountRequest()
		policyId         string
		orgSubAccountUin string
	)

	if v, ok := d.GetOkExists("policy_id"); ok {
		request.PolicyId = helper.IntInt64(v.(int))
		policyId = helper.IntToStr(v.(int))
	}

	if v, ok := d.GetOkExists("org_sub_account_uin"); ok {
		request.OrgSubAccountUins = append(request.OrgSubAccountUins, helper.IntInt64(v.(int)))
		orgSubAccountUin = helper.IntToStr(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().BindOrganizationPolicySubAccountWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create organization members auth policy attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{policyId, orgSubAccountUin}, tccommon.FILED_SP))
	return resourceTencentCloudOrganizationMemberAuthPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudOrganizationMemberAuthPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_member_auth_policy_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	policyId := idSplit[0]
	orgSubAccountUin := idSplit[1]

	respData, err := service.DescribeOrganizationMembersAuthPolicyAttachmentById(ctx, policyId, orgSubAccountUin)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_organization_member_auth_policy_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if len(respData) != 1 {
		return fmt.Errorf("Query organization members auth policy attachment by id return more than one.")
	}

	for _, item := range respData {
		if item.PolicyId != nil {
			_ = d.Set("policy_id", item.PolicyId)
		}

		if item.OrgSubAccountUin != nil {
			_ = d.Set("org_sub_account_uin", item.OrgSubAccountUin)
		}

		if item.IdentityId != nil {
			_ = d.Set("identity_id", item.IdentityId)
		}

		if item.IdentityRoleName != nil {
			_ = d.Set("identity_role_name", item.IdentityRoleName)
		}

		if item.IdentityRoleAliasName != nil {
			_ = d.Set("identity_role_alias_name", item.IdentityRoleAliasName)
		}

		if item.CreateTime != nil {
			_ = d.Set("create_time", item.CreateTime)
		}

		if item.PolicyName != nil {
			_ = d.Set("policy_name", item.PolicyName)
		}

		if item.MemberUin != nil {
			_ = d.Set("member_uin", item.MemberUin)
		}

		if item.MemberName != nil {
			_ = d.Set("member_name", item.MemberName)
		}

		if item.OrgSubAccountName != nil {
			_ = d.Set("org_sub_account_name", item.OrgSubAccountName)
		}

		if item.BindType != nil {
			_ = d.Set("bind_type", item.BindType)
		}
	}

	return nil
}

func resourceTencentCloudOrganizationMemberAuthPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_member_auth_policy_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = organizationv20210331.NewCancelOrganizationPolicySubAccountRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	policyId := idSplit[0]
	orgSubAccountUin := idSplit[1]

	request.PolicyId = helper.StrToInt64Point(policyId)
	request.OrgSubAccountUins = append(request.OrgSubAccountUins, helper.StrToInt64Point(orgSubAccountUin))
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().CancelOrganizationPolicySubAccountWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete organization members auth policy attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

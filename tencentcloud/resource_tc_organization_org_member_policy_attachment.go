package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudOrganizationOrgMemberPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgMemberPolicyAttachmentCreate,
		Read:   resourceTencentCloudOrganizationOrgMemberPolicyAttachmentRead,
		Delete: resourceTencentCloudOrganizationOrgMemberPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"member_uins": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Member Uin list. Up to 10.",
			},

			"policy_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Policy name.The maximum length is 128 characters, supporting English letters, numbers, and symbols +=,.@_-.",
			},

			"identity_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Organization identity ID.",
			},

			"description": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Notes.The maximum length is 128 characters.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgMemberPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_policy_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = organization.NewCreateOrganizationMembersPolicyRequest()
		response = organization.NewCreateOrganizationMembersPolicyResponse()
	)
	if v, ok := d.GetOk("member_uins"); ok {
		memberUinsSet := v.(*schema.Set).List()
		for i := range memberUinsSet {
			memberUins := memberUinsSet[i].(int)
			request.MemberUins = append(request.MemberUins, helper.IntInt64(memberUins))
		}
	}

	if v, ok := d.GetOk("policy_name"); ok {
		request.PolicyName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("identity_id"); ok {
		request.IdentityId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().CreateOrganizationMembersPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization orgMemberPolicyAttachment failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil || response.Response.PolicyId == nil {
		return fmt.Errorf("policy id is null")
	}
	policyId := *response.Response.PolicyId
	d.SetId(helper.Int64ToStr(policyId))

	return resourceTencentCloudOrganizationOrgMemberPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudOrganizationOrgMemberPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_policy_attachment.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudOrganizationOrgMemberPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_policy_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}
	orgMemberPolicyAttachmentId := d.Id()

	if err := service.DeleteOrganizationOrgMemberPolicyAttachmentById(ctx, orgMemberPolicyAttachmentId); err != nil {
		return err
	}

	return nil
}

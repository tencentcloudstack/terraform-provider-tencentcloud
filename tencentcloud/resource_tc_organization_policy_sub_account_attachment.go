/*
Provides a resource to create a organization policy_sub_account_attachment

Example Usage

```hcl
resource "tencentcloud_organization_policy_sub_account_attachment" "policy_sub_account_attachment" {
  policy_id = &lt;nil&gt;
  org_sub_account_uins = &lt;nil&gt;
  member_uin = &lt;nil&gt;
  org_sub_account_uin = &lt;nil&gt;
  policy_name = &lt;nil&gt;
  identity_id = &lt;nil&gt;
  identity_role_name = &lt;nil&gt;
  identity_role_alias_name = &lt;nil&gt;
  create_time = &lt;nil&gt;
  update_time = &lt;nil&gt;
  org_sub_account_name = &lt;nil&gt;
}
```

Import

organization policy_sub_account_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_organization_policy_sub_account_attachment.policy_sub_account_attachment policy_sub_account_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudOrganizationPolicySubAccountAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationPolicySubAccountAttachmentCreate,
		Read:   resourceTencentCloudOrganizationPolicySubAccountAttachmentRead,
		Delete: resourceTencentCloudOrganizationPolicySubAccountAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Policy ID.",
			},

			"org_sub_account_uins": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Organization administrator sub account uin list.",
			},

			"member_uin": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Organization member uin.",
			},

			"org_sub_account_uin": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Organization administrator sub account uin.",
			},

			"policy_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Policy name.",
			},

			"identity_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Manage Identity ID.",
			},

			"identity_role_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Identity role name.",
			},

			"identity_role_alias_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Identity role alias name.",
			},

			"create_time": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Creation time.",
			},

			"update_time": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},

			"org_sub_account_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Organization administrator sub account name.",
			},
		},
	}
}

func resourceTencentCloudOrganizationPolicySubAccountAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_policy_sub_account_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = organization.NewBindOrganizationMemberAuthAccountRequest()
		response = organization.NewBindOrganizationMemberAuthAccountResponse()
		policyId int
	)
	if v, ok := d.GetOkExists("policy_id"); ok {
		policyId = v.(int)
		request.PolicyId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("org_sub_account_uins"); ok {
		orgSubAccountUinsSet := v.(*schema.Set).List()
		for i := range orgSubAccountUinsSet {
			orgSubAccountUins := orgSubAccountUinsSet[i].(int)
			request.OrgSubAccountUins = append(request.OrgSubAccountUins, helper.IntUint64(orgSubAccountUins))
		}
	}

	if v, ok := d.GetOkExists("member_uin"); ok {
		request.MemberUin = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("org_sub_account_uin"); ok {
		request.OrgSubAccountUin = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("policy_name"); ok {
		request.PolicyName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("identity_id"); ok {
		request.IdentityId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("identity_role_name"); ok {
		request.IdentityRoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("identity_role_alias_name"); ok {
		request.IdentityRoleAliasName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("create_time"); ok {
		request.CreateTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("update_time"); ok {
		request.UpdateTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("org_sub_account_name"); ok {
		request.OrgSubAccountName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().BindOrganizationMemberAuthAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization policySubAccountAttachment failed, reason:%+v", logId, err)
		return err
	}

	policyId = *response.Response.PolicyId
	d.SetId(strings.Join([]string{helper.Int64ToStr(int64(policyId))}, FILED_SP))

	return resourceTencentCloudOrganizationPolicySubAccountAttachmentRead(d, meta)
}

func resourceTencentCloudOrganizationPolicySubAccountAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_policy_sub_account_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]
	orgSubAccountUins := idSplit[1]

	policySubAccountAttachment, err := service.DescribeOrganizationPolicySubAccountAttachmentById(ctx, policyId, orgSubAccountUins)
	if err != nil {
		return err
	}

	if policySubAccountAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationPolicySubAccountAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if policySubAccountAttachment.PolicyId != nil {
		_ = d.Set("policy_id", policySubAccountAttachment.PolicyId)
	}

	if policySubAccountAttachment.OrgSubAccountUins != nil {
		_ = d.Set("org_sub_account_uins", policySubAccountAttachment.OrgSubAccountUins)
	}

	if policySubAccountAttachment.MemberUin != nil {
		_ = d.Set("member_uin", policySubAccountAttachment.MemberUin)
	}

	if policySubAccountAttachment.OrgSubAccountUin != nil {
		_ = d.Set("org_sub_account_uin", policySubAccountAttachment.OrgSubAccountUin)
	}

	if policySubAccountAttachment.PolicyName != nil {
		_ = d.Set("policy_name", policySubAccountAttachment.PolicyName)
	}

	if policySubAccountAttachment.IdentityId != nil {
		_ = d.Set("identity_id", policySubAccountAttachment.IdentityId)
	}

	if policySubAccountAttachment.IdentityRoleName != nil {
		_ = d.Set("identity_role_name", policySubAccountAttachment.IdentityRoleName)
	}

	if policySubAccountAttachment.IdentityRoleAliasName != nil {
		_ = d.Set("identity_role_alias_name", policySubAccountAttachment.IdentityRoleAliasName)
	}

	if policySubAccountAttachment.CreateTime != nil {
		_ = d.Set("create_time", policySubAccountAttachment.CreateTime)
	}

	if policySubAccountAttachment.UpdateTime != nil {
		_ = d.Set("update_time", policySubAccountAttachment.UpdateTime)
	}

	if policySubAccountAttachment.OrgSubAccountName != nil {
		_ = d.Set("org_sub_account_name", policySubAccountAttachment.OrgSubAccountName)
	}

	return nil
}

func resourceTencentCloudOrganizationPolicySubAccountAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_policy_sub_account_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]
	orgSubAccountUins := idSplit[1]

	if err := service.DeleteOrganizationPolicySubAccountAttachmentById(ctx, policyId, orgSubAccountUins); err != nil {
		return err
	}

	return nil
}

/*
Provides a resource to create a organization org_member_policy

Example Usage

```hcl
resource "tencentcloud_organization_org_member_policy" "org_member_policy" {
  member_uin = &lt;nil&gt;
  policy_name = &lt;nil&gt;
  identity_id = &lt;nil&gt;
  description = &lt;nil&gt;
}
```

Import

organization org_member_policy can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_member_policy.org_member_policy org_member_policy_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudOrganizationOrgMemberPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgMemberPolicyCreate,
		Read:   resourceTencentCloudOrganizationOrgMemberPolicyRead,
		Update: resourceTencentCloudOrganizationOrgMemberPolicyUpdate,
		Delete: resourceTencentCloudOrganizationOrgMemberPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"member_uin": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Organization member uin.",
			},

			"policy_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Policy name.The maximum length is 128 characters, supporting English letters, numbers, and symbols +=,.@_-.",
			},

			"identity_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Organization identity ID.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Notes.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgMemberPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_policy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = organization.NewCreateOrganizationMemberPolicyRequest()
		response = organization.NewCreateOrganizationMemberPolicyResponse()
		policyId int
	)
	if v, ok := d.GetOkExists("member_uin"); ok {
		request.MemberUin = helper.IntInt64(v.(int))
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().CreateOrganizationMemberPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization orgMemberPolicy failed, reason:%+v", logId, err)
		return err
	}

	policyId = *response.Response.PolicyId
	d.SetId(helper.Int64ToStr(policyId))

	return resourceTencentCloudOrganizationOrgMemberPolicyRead(d, meta)
}

func resourceTencentCloudOrganizationOrgMemberPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	orgMemberPolicyId := d.Id()

	orgMemberPolicy, err := service.DescribeOrganizationOrgMemberPolicyById(ctx, policyId)
	if err != nil {
		return err
	}

	if orgMemberPolicy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgMemberPolicy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if orgMemberPolicy.MemberUin != nil {
		_ = d.Set("member_uin", orgMemberPolicy.MemberUin)
	}

	if orgMemberPolicy.PolicyName != nil {
		_ = d.Set("policy_name", orgMemberPolicy.PolicyName)
	}

	if orgMemberPolicy.IdentityId != nil {
		_ = d.Set("identity_id", orgMemberPolicy.IdentityId)
	}

	if orgMemberPolicy.Description != nil {
		_ = d.Set("description", orgMemberPolicy.Description)
	}

	return nil
}

func resourceTencentCloudOrganizationOrgMemberPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_policy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"member_uin", "policy_name", "identity_id", "description"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudOrganizationOrgMemberPolicyRead(d, meta)
}

func resourceTencentCloudOrganizationOrgMemberPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_policy.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}
	orgMemberPolicyId := d.Id()

	if err := service.DeleteOrganizationOrgMemberPolicyById(ctx, policyId); err != nil {
		return err
	}

	return nil
}

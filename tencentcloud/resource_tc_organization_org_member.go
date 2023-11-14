/*
Provides a resource to create a organization org_member

Example Usage

```hcl
resource "tencentcloud_organization_org_member" "org_member" {
  name = &lt;nil&gt;
  policy_type = "Financial"
  permission_ids =
  node_id =
  account_name = ""
  remark = ""
  record_id =
  pay_uin = ""
  identity_role_i_d =
  auth_relation_id =
}
```

Import

organization org_member can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_member.org_member org_member_id
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

func resourceTencentCloudOrganizationOrgMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgMemberCreate,
		Read:   resourceTencentCloudOrganizationOrgMemberRead,
		Update: resourceTencentCloudOrganizationOrgMemberUpdate,
		Delete: resourceTencentCloudOrganizationOrgMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Member name.The maximum length is 25 characters, supporting English letters, numbers, Chinese characters, and symbols +@,&amp;amp;amp;._[]-:,.",
			},

			"policy_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Organization policy type.- `Financial`: Financial management policy.",
			},

			"permission_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Financial management permission IDs.Valid values:- `1`: View bill.- `2`: Check balance.- `3`: Fund transfer.- `4`: Combine bill.- `5`: Issue an invoice.- `6`: Inherit discount.- `7`: Pay on behalf.- `8`: Analysis cost.value 1,2 is required.",
			},

			"node_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Organization node ID.",
			},

			"account_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Member account name.The maximum length is 25 characters, supporting English letters, numbers, Chinese characters, and symbols +@,&amp;amp;amp;._[]-:,.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Notes.",
			},

			"record_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Create member record ID.When create failed and needs to be recreated, is required.",
			},

			"pay_uin": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The uin which is payment account on behalf.When `PermissionIds` contains 7, is required.",
			},

			"identity_role_i_d": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Manage Identity IDs.",
			},

			"auth_relation_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Auth relationships Id.When creating members for different auth, it is necessary to.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgMemberCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = organization.NewCreateOrganizationMemberRequest()
		response = organization.NewCreateOrganizationMemberResponse()
		uin      int
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("policy_type"); ok {
		request.PolicyType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("permission_ids"); ok {
		permissionIdsSet := v.(*schema.Set).List()
		for i := range permissionIdsSet {
			permissionIds := permissionIdsSet[i].(int)
			request.PermissionIds = append(request.PermissionIds, helper.IntUint64(permissionIds))
		}
	}

	if v, ok := d.GetOkExists("node_id"); ok {
		request.NodeId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("account_name"); ok {
		request.AccountName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("record_id"); ok {
		request.RecordId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("pay_uin"); ok {
		request.PayUin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("identity_role_i_d"); ok {
		identityRoleIDSet := v.(*schema.Set).List()
		for i := range identityRoleIDSet {
			identityRoleID := identityRoleIDSet[i].(int)
			request.IdentityRoleID = append(request.IdentityRoleID, helper.IntUint64(identityRoleID))
		}
	}

	if v, ok := d.GetOkExists("auth_relation_id"); ok {
		request.AuthRelationId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().CreateOrganizationMember(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization orgMember failed, reason:%+v", logId, err)
		return err
	}

	uin = *response.Response.Uin
	d.SetId(helper.Int64ToStr(uin))

	return resourceTencentCloudOrganizationOrgMemberRead(d, meta)
}

func resourceTencentCloudOrganizationOrgMemberRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	orgMemberId := d.Id()

	orgMember, err := service.DescribeOrganizationOrgMemberById(ctx, uin)
	if err != nil {
		return err
	}

	if orgMember == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgMember` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if orgMember.Name != nil {
		_ = d.Set("name", orgMember.Name)
	}

	if orgMember.PolicyType != nil {
		_ = d.Set("policy_type", orgMember.PolicyType)
	}

	if orgMember.PermissionIds != nil {
		_ = d.Set("permission_ids", orgMember.PermissionIds)
	}

	if orgMember.NodeId != nil {
		_ = d.Set("node_id", orgMember.NodeId)
	}

	if orgMember.AccountName != nil {
		_ = d.Set("account_name", orgMember.AccountName)
	}

	if orgMember.Remark != nil {
		_ = d.Set("remark", orgMember.Remark)
	}

	if orgMember.RecordId != nil {
		_ = d.Set("record_id", orgMember.RecordId)
	}

	if orgMember.PayUin != nil {
		_ = d.Set("pay_uin", orgMember.PayUin)
	}

	if orgMember.IdentityRoleID != nil {
		_ = d.Set("identity_role_i_d", orgMember.IdentityRoleID)
	}

	if orgMember.AuthRelationId != nil {
		_ = d.Set("auth_relation_id", orgMember.AuthRelationId)
	}

	return nil
}

func resourceTencentCloudOrganizationOrgMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		moveOrganizationNodeMembersRequest  = organization.NewMoveOrganizationNodeMembersRequest()
		moveOrganizationNodeMembersResponse = organization.NewMoveOrganizationNodeMembersResponse()
	)

	orgMemberId := d.Id()

	request.Uin = &uin

	immutableArgs := []string{"name", "policy_type", "permission_ids", "node_id", "account_name", "remark", "record_id", "pay_uin", "identity_role_i_d", "auth_relation_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("policy_type") {
		if v, ok := d.GetOk("policy_type"); ok {
			request.PolicyType = helper.String(v.(string))
		}
	}

	if d.HasChange("permission_ids") {
		if v, ok := d.GetOk("permission_ids"); ok {
			permissionIdsSet := v.(*schema.Set).List()
			for i := range permissionIdsSet {
				permissionIds := permissionIdsSet[i].(int)
				request.PermissionIds = append(request.PermissionIds, helper.IntUint64(permissionIds))
			}
		}
	}

	if d.HasChange("node_id") {
		if v, ok := d.GetOkExists("node_id"); ok {
			request.NodeId = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	if d.HasChange("pay_uin") {
		if v, ok := d.GetOk("pay_uin"); ok {
			request.PayUin = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().MoveOrganizationNodeMembers(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update organization orgMember failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudOrganizationOrgMemberRead(d, meta)
}

func resourceTencentCloudOrganizationOrgMemberDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}
	orgMemberId := d.Id()

	if err := service.DeleteOrganizationOrgMemberById(ctx, uin); err != nil {
		return err
	}

	return nil
}

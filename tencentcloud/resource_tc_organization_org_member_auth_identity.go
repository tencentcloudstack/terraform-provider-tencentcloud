/*
Provides a resource to create a organization org_member_auth_identity

Example Usage

```hcl
resource "tencentcloud_organization_org_member_auth_identity" "org_member_auth_identity" {
  member_uins = &lt;nil&gt;
  identity_ids = &lt;nil&gt;
}
```

Import

organization org_member_auth_identity can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_member_auth_identity.org_member_auth_identity org_member_auth_identity_id
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

func resourceTencentCloudOrganizationOrgMemberAuthIdentity() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgMemberAuthIdentityCreate,
		Read:   resourceTencentCloudOrganizationOrgMemberAuthIdentityRead,
		Delete: resourceTencentCloudOrganizationOrgMemberAuthIdentityDelete,
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

func resourceTencentCloudOrganizationOrgMemberAuthIdentityCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_auth_identity.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = organization.NewCreateOrganizationMemberAuthIdentityRequest()
		response   = organization.NewCreateOrganizationMemberAuthIdentityResponse()
		memberUin  int
		identityId int
	)
	if v, ok := d.GetOk("member_uins"); ok {
		memberUinsSet := v.(*schema.Set).List()
		for i := range memberUinsSet {
			memberUins := memberUinsSet[i].(int)
			request.MemberUins = append(request.MemberUins, helper.IntUint64(memberUins))
		}
	}

	if v, ok := d.GetOk("identity_ids"); ok {
		identityIdsSet := v.(*schema.Set).List()
		for i := range identityIdsSet {
			identityIds := identityIdsSet[i].(int)
			request.IdentityIds = append(request.IdentityIds, helper.IntUint64(identityIds))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().CreateOrganizationMemberAuthIdentity(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization orgMemberAuthIdentity failed, reason:%+v", logId, err)
		return err
	}

	memberUin = *response.Response.MemberUin
	d.SetId(strings.Join([]string{helper.Int64ToStr(int64(memberUin)), helper.Int64ToStr(int64(identityId))}, FILED_SP))

	return resourceTencentCloudOrganizationOrgMemberAuthIdentityRead(d, meta)
}

func resourceTencentCloudOrganizationOrgMemberAuthIdentityRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_auth_identity.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	memberUin := idSplit[0]
	identityId := idSplit[1]

	orgMemberAuthIdentity, err := service.DescribeOrganizationOrgMemberAuthIdentityById(ctx, memberUin, identityId)
	if err != nil {
		return err
	}

	if orgMemberAuthIdentity == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgMemberAuthIdentity` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if orgMemberAuthIdentity.MemberUins != nil {
		_ = d.Set("member_uins", orgMemberAuthIdentity.MemberUins)
	}

	if orgMemberAuthIdentity.IdentityIds != nil {
		_ = d.Set("identity_ids", orgMemberAuthIdentity.IdentityIds)
	}

	return nil
}

func resourceTencentCloudOrganizationOrgMemberAuthIdentityDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_auth_identity.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	memberUin := idSplit[0]
	identityId := idSplit[1]

	if err := service.DeleteOrganizationOrgMemberAuthIdentityById(ctx, memberUin, identityId); err != nil {
		return err
	}

	return nil
}

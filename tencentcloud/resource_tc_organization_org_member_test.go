package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudOrganizationOrgMemberResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgMember,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member.org_member", "id")),
			},
			{
				ResourceName:      "tencentcloud_organization_org_member.org_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgMember = `

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

`

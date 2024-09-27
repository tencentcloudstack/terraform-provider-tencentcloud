package tco_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationOrgMemberResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgMember,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member.org_member", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member.org_member", "name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member.org_member", "node_id", "2023047"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member.org_member", "permission_ids.#"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member.org_member", "policy_type", "Financial"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member.org_member", "remark", "for terraform test"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member.org_member", "tags.tf-test1", "org_member1"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member.org_member", "tags.tf-test2", "org_member2"),
				),
			},
			{
				Config: testAccOrganizationOrgMemberUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member.org_member", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member.org_member", "name", "tf_example_1"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member.org_member", "node_id", "2023047"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_member.org_member", "permission_ids.#"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member.org_member", "policy_type", "Financial"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member.org_member", "remark", "for terraform test"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_member.org_member", "tags.tf-test1", "org_member_update"),
				),
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
  name            = "tf_example"
  node_id         = 2023047
  permission_ids  = [
    1,
    2,
    4,
	5,
  ]
  policy_type     = "Financial"
  remark          = "for terraform test"
  tags = {
    "tf-test1" = "org_member1"
	"tf-test2" = "org_member2"
  }
}
`
const testAccOrganizationOrgMemberUpdate = `

resource "tencentcloud_organization_org_member" "org_member" {
  name            = "tf_example_1"
  node_id         = 2023047
  permission_ids  = [
    1,
    2,
    4,
  ]
  policy_type     = "Financial"
  remark          = "for terraform test"
  tags = {
    "tf-test1" = "org_member_update"
  }
}

`

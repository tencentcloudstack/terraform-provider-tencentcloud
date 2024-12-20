package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudAcceptJoinShareUnitInvitationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAcceptJoinShareUnitInvitationOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_accept_join_share_unit_invitation_operation.accept_join_share_unit_invitation_operation", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_accept_join_share_unit_invitation_operation.accept_join_share_unit_invitation_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAcceptJoinShareUnitInvitationOperation = `
resource "tencentcloud_organization_org_share_unit" "org_share_unit" {
  name = "iac-test-1"
  area = "ap-guangzhou"
  description = "iac-test"
}
resource "tencentcloud_accept_join_share_unit_invitation_operation" "accept_join_share_unit_invitation_operation" {
  unit_id = split("#", tencentcloud_organization_org_share_unit.org_share_unit.id)[1]
}
`

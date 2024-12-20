package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudRejectJoinShareUnitInvitationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccRejectJoinShareUnitInvitationOperation,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_reject_join_share_unit_invitation_operation.reject_join_share_unit_invitation_operation", "id")),
		}, {
			ResourceName:      "tencentcloud_reject_join_share_unit_invitation_operation.reject_join_share_unit_invitation_operation",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccRejectJoinShareUnitInvitationOperation = `
resource "tencentcloud_organization_org_share_unit" "org_share_unit" {
  name = "iac-test-1"
  area = "ap-guangzhou"
  description = "iac-test"
}
resource "tencentcloud_reject_join_share_unit_invitation_operation" "reject_join_share_unit_invitation_operation" {
  unit_id = split("#", tencentcloud_organization_org_share_unit.org_share_unit.id)[1]
}
`

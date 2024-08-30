package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterUserSyncProvisioningResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCenterUserSyncProvisioning,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "deletion_strategy", "Keep"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "description", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "duplication_strategy", "TakeOver"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "principal_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "principal_name"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "principal_type", "User"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "status", "Enabled"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "target_name"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "target_type", "MemberUin"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "target_uin"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "update_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "user_provisioning_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "zone_id"),
				),
			},
			{
				Config: testAccIdentityCenterUserSyncProvisioningUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "deletion_strategy", "Keep"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "description", "tf-test-update"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "duplication_strategy", "TakeOver"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "principal_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "principal_name"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "principal_type", "User"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "status", "Enabled"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "target_name"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "target_type", "MemberUin"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "target_uin"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "update_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "user_provisioning_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning", "zone_id"),
				),
			},
			{
				ResourceName: "tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning",
				ImportState:  true,
			},
		},
	})
}

const testAccIdentityCenterUserSyncProvisioning = `
resource "tencentcloud_identity_center_user" "identity_center_user" {
    zone_id = "z-s64jh54hbcra"
    user_name = "sync-user"
    description = "test"
}
data "tencentcloud_organization_members" "members" {}

resource "tencentcloud_identity_center_user_sync_provisioning" "identity_center_user_sync_provisioning" {
  zone_id = "z-s64jh54hbcra"
  description = "tf-test"
  deletion_strategy = "Keep"
  duplication_strategy = "TakeOver"
  principal_id = tencentcloud_identity_center_user.identity_center_user.user_id
  principal_type = "User"
  target_uin = data.tencentcloud_organization_members.members.items.0.member_uin
  target_type = "MemberUin"
}
`

const testAccIdentityCenterUserSyncProvisioningUpdate = `
resource "tencentcloud_identity_center_user" "identity_center_user" {
    zone_id = "z-s64jh54hbcra"
    user_name = "sync-user"
    description = "test"
}
data "tencentcloud_organization_members" "members" {}

resource "tencentcloud_identity_center_user_sync_provisioning" "identity_center_user_sync_provisioning" {
  zone_id = "z-s64jh54hbcra"
  description = "tf-test-update"
  deletion_strategy = "Keep"
  duplication_strategy = "TakeOver"
  principal_id = tencentcloud_identity_center_user.identity_center_user.user_id
  principal_type = "User"
  target_uin = data.tencentcloud_organization_members.members.items.0.member_uin
  target_type = "MemberUin"
}
`

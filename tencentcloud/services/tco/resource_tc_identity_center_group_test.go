package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCenterGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_group.identity_center_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_group.identity_center_group", "zone_id", "z-s64jh54hbcra"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_group.identity_center_group", "group_name", "tf-test-group"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_group.identity_center_group", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_group.identity_center_group", "group_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_group.identity_center_group", "update_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_group.identity_center_group", "group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_group.identity_center_group", "member_count"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_group.identity_center_group", "description", "test"),
				),
			},
			{
				Config: testAccIdentityCenterGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_group.identity_center_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_group.identity_center_group", "description", "test_update"),
				),
			},
			{
				ResourceName:      "tencentcloud_identity_center_group.identity_center_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIdentityCenterGroup = `
resource "tencentcloud_identity_center_group" "identity_center_group" {
	zone_id = "z-s64jh54hbcra"
    group_name = "tf-test-group"
    description = "test"
}
`

const testAccIdentityCenterGroupUpdate = `
resource "tencentcloud_identity_center_group" "identity_center_group" {
	zone_id = "z-s64jh54hbcra"
    group_name = "tf-test-group"
    description = "test_update"
}
`

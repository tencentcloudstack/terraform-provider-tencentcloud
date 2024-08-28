package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterRoleConfigurationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCenterRoleConfiguration,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_configuration.identity_center_role_configuration", "id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_role_configuration.identity_center_role_configuration", "role_configuration_name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_role_configuration.identity_center_role_configuration", "description", "test"),
				),
			},
			{
				Config: testAccIdentityCenterRoleConfigurationUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_role_configuration.identity_center_role_configuration", "id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_role_configuration.identity_center_role_configuration", "role_configuration_name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_role_configuration.identity_center_role_configuration", "description", "test-update"),
				),
			},
			{
				ResourceName:      "tencentcloud_identity_center_role_configuration.identity_center_role_configuration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIdentityCenterRoleConfiguration = `
resource "tencentcloud_identity_center_role_configuration" "identity_center_role_configuration" {
    zone_id = "z-s64jh54hbcra"
    role_configuration_name = "tf-test"
    description = "test"
}
`

const testAccIdentityCenterRoleConfigurationUpdate = `
resource "tencentcloud_identity_center_role_configuration" "identity_center_role_configuration" {
    zone_id = "z-s64jh54hbcra"
    role_configuration_name = "tf-test"
    description = "test-update"
}
`

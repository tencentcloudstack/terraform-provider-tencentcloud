package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudRoleConfigurationProvisioningsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleConfigurationProvisioningsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_role_configuration_provisionings.role_configuration_provisionings"),
					resource.TestCheckResourceAttr("data.tencentcloud_role_configuration_provisionings.role_configuration_provisionings", "role_configuration_provisionings.#", "1"),
				),
			},
		},
	})
}

const testAccRoleConfigurationProvisioningsDataSource = testAccProvisionRoleConfigurationOperation + `
data "tencentcloud_role_configuration_provisionings" "role_configuration_provisionings" {
  zone_id               = "z-s64jh54hbcra"
  role_configuration_id = tencentcloud_identity_center_role_configuration.identity_center_role_configuration.role_configuration_id
  depends_on            = [tencentcloud_provision_role_configuration_operation.provision_role_configuration_operation]
}
`

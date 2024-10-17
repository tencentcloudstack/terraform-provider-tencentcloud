package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterRoleConfigurationsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccIdentityCenterRoleConfigurationsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_identity_center_role_configurations.identity_center_role_configurations"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_identity_center_role_configurations.identity_center_role_configurations", "role_configurations.#"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_identity_center_role_configurations.identity_center_role_configurations", "role_configurations.0.role_configuration_id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_identity_center_role_configurations.identity_center_role_configurations", "role_configurations.0.role_configuration_name"),
			),
		}},
	})
}

const testAccIdentityCenterRoleConfigurationsDataSource = `
data "tencentcloud_identity_center_role_configurations" "identity_center_role_configurations" {
    zone_id = "z-s64jh54hbcra"
}
`

package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccIdentityCenterGroupsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_identity_center_groups.identity_center_groups"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_identity_center_groups.identity_center_groups", "groups.#"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_identity_center_groups.identity_center_groups", "groups.0.group_id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_identity_center_groups.identity_center_groups", "groups.0.group_name"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_identity_center_groups.identity_center_groups", "groups.0.group_type"),
			),
		}},
	})
}

const testAccIdentityCenterGroupsDataSource = `
data "tencentcloud_identity_center_groups" "identity_center_groups" {
    zone_id = "z-s64jh54hbcra"
}
`

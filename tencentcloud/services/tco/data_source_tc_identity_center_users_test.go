package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterUsersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccIdentityCenterUsersDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_identity_center_users.identity_center_users"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_identity_center_users.identity_center_users", "users.0.user_id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_identity_center_users.identity_center_users", "users.0.user_name"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_identity_center_users.identity_center_users", "users.0.user_status"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_identity_center_users.identity_center_users", "users.0.user_type"),
			),
		}},
	})
}

const testAccIdentityCenterUsersDataSource = `
data "tencentcloud_identity_center_users" "identity_center_users" {
    zone_id = "z-s64jh54hbcra"
}
`

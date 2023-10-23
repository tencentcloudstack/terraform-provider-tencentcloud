package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeUserRolesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeUserRolesDataSource,
				Check: resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_user_roles.describe_user_roles"),
					resource.TestCheckResourceAttr("data.tencentcloud_dlc_describe_user_roles.describe_user_roles", "fuzzy", "1")),
			},
		},
	})
}

const testAccDlcDescribeUserRolesDataSource = `

data "tencentcloud_dlc_describe_user_roles" "describe_user_roles" {
  fuzzy = "1"
  }

`

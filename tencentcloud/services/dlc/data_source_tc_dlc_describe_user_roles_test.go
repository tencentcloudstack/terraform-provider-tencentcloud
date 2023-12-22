package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeUserRolesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeUserRolesDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_user_roles.describe_user_roles"),
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

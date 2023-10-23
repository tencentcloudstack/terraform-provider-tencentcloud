package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeUserInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeUserInfoDataSource,
				Check: resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_user_info.describe_user_info"),
					resource.TestCheckResourceAttr("data.tencentcloud_dlc_describe_user_info.describe_user_info", "user_id", "100032772113"),
					resource.TestCheckResourceAttr("data.tencentcloud_dlc_describe_user_info.describe_user_info", "type", "Group"),
					resource.TestCheckResourceAttr("data.tencentcloud_dlc_describe_user_info.describe_user_info", "sort_by", "create-time"),
					resource.TestCheckResourceAttr("data.tencentcloud_dlc_describe_user_info.describe_user_info", "sorting", "desc"),
				),
			},
		},
	})
}

const testAccDlcDescribeUserInfoDataSource = `

data "tencentcloud_dlc_describe_user_info" "describe_user_info" {
  user_id = "100032772113"
  type = "Group"
  sort_by = "create-time"
  sorting = "desc"
}
`

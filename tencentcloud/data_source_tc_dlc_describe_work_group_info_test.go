package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeWorkGroupInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeWorkGroupInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_work_group_info.describe_work_group_info")),
			},
		},
	})
}

const testAccDlcDescribeWorkGroupInfoDataSource = `

data "tencentcloud_dlc_describe_work_group_info" "describe_work_group_info" {
  work_group_id = 23181
  type = "User"
  sort_by = "create-time"
  sorting = "desc"
  }

`

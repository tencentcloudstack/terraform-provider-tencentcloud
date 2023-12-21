package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeWorkGroupInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeWorkGroupInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_work_group_info.describe_work_group_info")),
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

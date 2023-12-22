package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeEngineUsageInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeEngineUsageInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_engine_usage_info.describe_engine_usage_info")),
			},
		},
	})
}

const testAccDlcDescribeEngineUsageInfoDataSource = `

data "tencentcloud_dlc_describe_engine_usage_info" "describe_engine_usage_info" {
  data_engine_id = "DataEngine-cgkvbas6"
    }

`

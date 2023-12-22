package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeDataEngineDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeDataEngineDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_data_engine.describe_data_engine"),
					resource.TestCheckResourceAttr("data.tencentcloud_dlc_describe_data_engine.describe_data_engine", "data_engine_name", "iac-test-spark"),
				),
			},
		},
	})
}

const testAccDlcDescribeDataEngineDataSource = `

data "tencentcloud_dlc_describe_data_engine" "describe_data_engine" {
  data_engine_name = "iac-test-spark"
  }

`

package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDlcDescribeDataEngineDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeDataEngineDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_data_engine.describe_data_engine")),
			},
		},
	})
}

const testAccDlcDescribeDataEngineDataSource = `

data "tencentcloud_dlc_describe_data_engine" "describe_data_engine" {
  data_engine_name = "iac-test-spark"
  }

`

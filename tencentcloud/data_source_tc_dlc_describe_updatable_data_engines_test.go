package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeUpdatableDataEnginesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeUpdatableDataEnginesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_updatable_data_engines.describe_updatable_data_engines")),
			},
		},
	})
}

const testAccDlcDescribeUpdatableDataEnginesDataSource = `

data "tencentcloud_dlc_describe_updatable_data_engines" "describe_updatable_data_engines" {
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
  }
`

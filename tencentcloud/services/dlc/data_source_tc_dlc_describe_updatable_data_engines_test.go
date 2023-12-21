package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeUpdatableDataEnginesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeUpdatableDataEnginesDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_updatable_data_engines.describe_updatable_data_engines")),
			},
		},
	})
}

const testAccDlcDescribeUpdatableDataEnginesDataSource = `

data "tencentcloud_dlc_describe_updatable_data_engines" "describe_updatable_data_engines" {
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
  }
`

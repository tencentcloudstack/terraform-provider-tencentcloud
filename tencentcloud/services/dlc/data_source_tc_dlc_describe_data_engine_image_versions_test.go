package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeDataEngineImageVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeDataEngineImageVersionsDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_data_engine_image_versions.describe_data_engine_image_versions"),
					resource.TestCheckResourceAttr("data.tencentcloud_dlc_describe_data_engine_image_versions.describe_data_engine_image_versions", "engine_type", "SparkBatch"),
				),
			},
		},
	})
}

const testAccDlcDescribeDataEngineImageVersionsDataSource = `

data "tencentcloud_dlc_describe_data_engine_image_versions" "describe_data_engine_image_versions" {
  engine_type = "SparkBatch"
  }

`

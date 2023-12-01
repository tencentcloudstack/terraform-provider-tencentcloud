package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeDataEnginePythonSparkImagesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeDataEnginePythonSparkImagesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_data_engine_python_spark_images.describe_data_engine_python_spark_images")),
			},
		},
	})
}

const testAccDlcDescribeDataEnginePythonSparkImagesDataSource = `

data "tencentcloud_dlc_describe_data_engine_python_spark_images" "describe_data_engine_python_spark_images" {
  child_image_version_id = "f54fba71-5f9c-4dfe-a565-004d7b6d3864"
  }

`

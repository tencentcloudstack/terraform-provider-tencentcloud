package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcStandardEngineResourceGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcStandardEngineResourceGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "engine_resource_group_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "data_engine_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "auto_launch"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "auto_pause"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "auto_pause_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "max_concurrency"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "resource_group_scene"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "spark_spec_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "spark_size"),
				),
			},
			{
				Config: testAccDlcStandardEngineResourceGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "engine_resource_group_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "data_engine_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "auto_launch"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "auto_pause"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "auto_pause_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "max_concurrency"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "resource_group_scene"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "spark_spec_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.example", "spark_size"),
				),
			},
		},
	})
}

const testAccDlcStandardEngineResourceGroup = `
resource "tencentcloud_dlc_standard_engine_resource_group" "example" {
  engine_resource_group_name = "tf-example"
  data_engine_name           = "tf-engine"
  auto_launch                = 0
  auto_pause                 = 0
  auto_pause_time            = 10
  static_config_pairs {
    config_item  = "key"
    config_value = "value"
  }

  dynamic_config_pairs {
    config_item  = "key"
    config_value = "value"
  }
  max_concurrency      = 5
  resource_group_scene = "SparkSQL"
  spark_spec_mode      = "fast"
  spark_size           = 16
}
`

const testAccDlcStandardEngineResourceGroupUpdate = `
resource "tencentcloud_dlc_standard_engine_resource_group" "example" {
  engine_resource_group_name = "tf-example"
  data_engine_name           = "tf-engine"
  auto_launch                = 0
  auto_pause                 = 0
  auto_pause_time            = 20
  static_config_pairs {
    config_item  = "key"
    config_value = "value"
  }

  dynamic_config_pairs {
    config_item  = "key"
    config_value = "value"
  }
  max_concurrency      = 10
  resource_group_scene = "SparkSQL"
  spark_spec_mode      = "fast"
  spark_size           = 16
}
`

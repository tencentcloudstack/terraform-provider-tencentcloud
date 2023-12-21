package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcSwitchDataEngineImageOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcSwitchDataEngineImageOperation,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_switch_data_engine_image_operation.switch_data_engine_image_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_switch_data_engine_image_operation.switch_data_engine_image_operation", "data_engine_id", "DataEngine-cgkvbas6"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_switch_data_engine_image_operation.switch_data_engine_image_operation", "new_image_version_id"),
				),
			},
			{
				Config: testAccDlcSwitchDataEngineImageOperationRecover,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_switch_data_engine_image_operation.switch_data_engine_image_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_switch_data_engine_image_operation.switch_data_engine_image_operation", "data_engine_id", "DataEngine-cgkvbas6"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_switch_data_engine_image_operation.switch_data_engine_image_operation", "new_image_version_id"),
				),
			},
		},
	})
}

const testAccDlcSwitchDataEngineImageOperation = `

data "tencentcloud_dlc_describe_data_engine_image_versions" "describe_data_engine_image_versions" {
  engine_type = "SparkBatch"
  }

resource "tencentcloud_dlc_switch_data_engine_image_operation" "switch_data_engine_image_operation" {
  data_engine_id = "DataEngine-cgkvbas6"
  new_image_version_id = data.tencentcloud_dlc_describe_data_engine_image_versions.describe_data_engine_image_versions.image_parent_versions.0.image_version_id
}

`
const testAccDlcSwitchDataEngineImageOperationRecover = `

data "tencentcloud_dlc_describe_data_engine_image_versions" "describe_data_engine_image_versions" {
  engine_type = "SparkBatch"
  }

resource "tencentcloud_dlc_switch_data_engine_image_operation" "switch_data_engine_image_operation" {
  data_engine_id = "DataEngine-cgkvbas6"
  new_image_version_id = data.tencentcloud_dlc_describe_data_engine_image_versions.describe_data_engine_image_versions.image_parent_versions.1.image_version_id
}

`

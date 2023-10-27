package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcModifyDataEngineDescriptionOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcModifyDataEngineDescriptionOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_modify_data_engine_description_operation.modify_data_engine_description_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_modify_data_engine_description_operation.modify_data_engine_description_operation", "data_engine_name", "iac-test-spark"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_modify_data_engine_description_operation.modify_data_engine_description_operation", "message", "test")),
			},
		},
	})
}

const testAccDlcModifyDataEngineDescriptionOperation = `

resource "tencentcloud_dlc_modify_data_engine_description_operation" "modify_data_engine_description_operation" {
  data_engine_name = "iac-test-spark"
  message = "test"
}

`

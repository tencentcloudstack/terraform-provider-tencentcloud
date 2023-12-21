package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcRollbackDataEngineImageOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcRollbackDataEngineImage,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_rollback_data_engine_image_operation.rollback_data_engine_image", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_rollback_data_engine_image_operation.rollback_data_engine_image", "data_engine_id", "DataEngine-cgkvbas6"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_rollback_data_engine_image_operation.rollback_data_engine_image", "from_record_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_rollback_data_engine_image_operation.rollback_data_engine_image", "to_record_id")),
			},
		},
	})
}

const testAccDlcRollbackDataEngineImage = `
data "tencentcloud_dlc_check_data_engine_image_can_be_rollback" "check_data_engine_image_can_be_rollback" {
  data_engine_id = "DataEngine-cgkvbas6"
}
resource "tencentcloud_dlc_rollback_data_engine_image_operation" "rollback_data_engine_image" {
  data_engine_id = "DataEngine-cgkvbas6"
  from_record_id = data.tencentcloud_dlc_check_data_engine_image_can_be_rollback.check_data_engine_image_can_be_rollback.from_record_id
  to_record_id = data.tencentcloud_dlc_check_data_engine_image_can_be_rollback.check_data_engine_image_can_be_rollback.to_record_id
}
`

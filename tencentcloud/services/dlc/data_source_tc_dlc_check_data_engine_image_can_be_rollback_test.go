package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcCheckDataEngineImageCanBeRollbackDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcCheckDataEngineImageCanBeRollbackDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_check_data_engine_image_can_be_rollback.check_data_engine_image_can_be_rollback"),
					resource.TestCheckResourceAttr("data.tencentcloud_dlc_check_data_engine_image_can_be_rollback.check_data_engine_image_can_be_rollback", "data_engine_id", "DataEngine-public-1308919341")),
			},
		},
	})
}

const testAccDlcCheckDataEngineImageCanBeRollbackDataSource = `

data "tencentcloud_dlc_check_data_engine_image_can_be_rollback" "check_data_engine_image_can_be_rollback" {
  data_engine_id = "DataEngine-public-1308919341"
      }

`

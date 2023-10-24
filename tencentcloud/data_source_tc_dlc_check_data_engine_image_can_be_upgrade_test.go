package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcCheckDataEngineImageCanBeUpgradeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcCheckDataEngineImageCanBeUpgradeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_check_data_engine_image_can_be_upgrade.check_data_engine_image_can_be_upgrade")),
			},
		},
	})
}

const testAccDlcCheckDataEngineImageCanBeUpgradeDataSource = `

data "tencentcloud_dlc_check_data_engine_image_can_be_upgrade" "check_data_engine_image_can_be_upgrade" {
  data_engine_id = "DataEngine-cgkvbas6"
    }

`

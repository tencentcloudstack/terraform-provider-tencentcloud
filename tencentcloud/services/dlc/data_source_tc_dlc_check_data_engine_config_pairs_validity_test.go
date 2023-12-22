package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcCheckDataEngineConfigPairsValidityDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcCheckDataEngineConfigPairsValidityDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_check_data_engine_config_pairs_validity.check_data_engine_config_pairs_validity"),
					resource.TestCheckResourceAttr("data.tencentcloud_dlc_check_data_engine_config_pairs_validity.check_data_engine_config_pairs_validity", "child_image_version_id", "f54fba71-5f9c-4dfe-a565-004d7b6d3864")),
			},
		},
	})
}

const testAccDlcCheckDataEngineConfigPairsValidityDataSource = `

data "tencentcloud_dlc_check_data_engine_config_pairs_validity" "check_data_engine_config_pairs_validity" {
  child_image_version_id = "f54fba71-5f9c-4dfe-a565-004d7b6d3864"
}
`

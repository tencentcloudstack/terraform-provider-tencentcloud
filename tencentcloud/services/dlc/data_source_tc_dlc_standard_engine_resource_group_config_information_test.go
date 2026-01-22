package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcStandardEngineResourceGroupConfigInformationDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccDlcStandardEngineResourceGroupConfigInformationDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_standard_engine_resource_group_config_information.example"),
			),
		}},
	})
}

const testAccDlcStandardEngineResourceGroupConfigInformationDataSource = `
data "tencentcloud_dlc_standard_engine_resource_group_config_information" "example" {
  sort_by = "create-time"
  sorting = "desc"
  filters {
    name = "engine-id"
    values = [
      "DataEngine-5plqp7q7"
    ]
  }
}
`

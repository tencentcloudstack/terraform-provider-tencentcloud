package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcDataEngineSessionParametersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccDlcDataEngineSessionParametersDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_data_engine_session_parameters.dlc_data_engine_session_parameters"),
			),
		}},
	})
}

const testAccDlcDataEngineSessionParametersDataSource = `
data "tencentcloud_dlc_data_engine_session_parameters" "example" {
  data_engine_id = "DataEngine-public-1308726196"
}
`

package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcDataEngineNetworkDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccDlcDataEngineNetworkDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_data_engine_network.example"),
			),
		}},
	})
}

const testAccDlcDataEngineNetworkDataSource = `
data "tencentcloud_dlc_data_engine_network" "example" {
  sort_by = "create-time"
  sorting = "desc"
  filters {
    name   = "engine-network-id"
    values = ["DataEngine_Network-g1sxyw8v"]
  }
}
`

package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcEngineNodeSpecificationsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccDlcEngineNodeSpecificationsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_engine_node_specifications.example"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_dlc_engine_node_specifications.example", "data_engine_name"),
			),
		}},
	})
}

const testAccDlcEngineNodeSpecificationsDataSource = `
data "tencentcloud_dlc_engine_node_specifications" "example" {
  data_engine_name = "tf-example"
}
`

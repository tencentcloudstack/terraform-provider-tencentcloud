package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcTaskResultDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccDlcTaskResultDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_task_result.example"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_dlc_task_result.example", "task_id"),
			),
		}},
	})
}

const testAccDlcTaskResultDataSource = `
data "tencentcloud_dlc_task_result" "example" {
  task_id = "fdd9c5fa21ca11eca6fb5254006c64af"
}
`

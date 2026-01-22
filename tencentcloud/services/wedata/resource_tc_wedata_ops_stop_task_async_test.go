package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsStopTaskAsyncResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataOpsStopTaskAsync,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_ops_stop_task_async.wedata_stop_ops_task_async", "id")),
		}, {
			ResourceName:      "tencentcloud_wedata_ops_stop_task_async.wedata_stop_ops_task_async",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccWedataOpsStopTaskAsync = `

resource "tencentcloud_wedata_ops_stop_task_async" "wedata_stop_ops_task_async" {
}
`

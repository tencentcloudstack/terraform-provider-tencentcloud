package wedata_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"testing"
)

func TestAccTencentCloudWedataTaskKillInstanceAsyncResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataTaskKillInstanceAsync,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_task_kill_instance_async.wedata_task_kill_instance_async", "id")),
		}, {
			ResourceName:      "tencentcloud_wedata_task_kill_instance_async.wedata_task_kill_instance_async",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccWedataTaskKillInstanceAsync = `

resource "tencentcloud_wedata_task_kill_instance_async" "wedata_task_kill_instance_async" {
}
`

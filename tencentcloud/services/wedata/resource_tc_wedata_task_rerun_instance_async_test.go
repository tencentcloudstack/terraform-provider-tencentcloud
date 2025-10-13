package wedata_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"testing"
)

func TestAccTencentCloudWedataTaskRerunInstanceAsyncResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataTaskRerunInstanceAsync,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_task_rerun_instance_async.wedata_task_rerun_instance_async", "id")),
		}, {
			ResourceName:      "tencentcloud_wedata_task_rerun_instance_async.wedata_task_rerun_instance_async",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccWedataTaskRerunInstanceAsync = `

resource "tencentcloud_wedata_task_rerun_instance_async" "wedata_task_rerun_instance_async" {
  redefine_param_list = {
  }
}
`

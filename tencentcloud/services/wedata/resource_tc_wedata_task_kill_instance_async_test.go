package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataTaskKillInstanceAsyncResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataTaskKillInstanceAsync,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_task_kill_instance_async.wedata_task_kill_instance_async", "id")),
			},
		},
	})
}

const testAccWedataTaskKillInstanceAsync = `

resource "tencentcloud_wedata_task_kill_instance_async" "wedata_task_kill_instance_async" {
  project_id        = "1859317240494305280"
  instance_key_list = ["20251013154418424_2025-10-13 18:10:00"]
}
`

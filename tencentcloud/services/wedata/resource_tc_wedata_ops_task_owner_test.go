package wedata_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"testing"
)

func TestAccTencentCloudWedataOpsTaskOwnerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataOpsTaskOwner,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner", "id")),
		}, {
			ResourceName:      "tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccWedataOpsTaskOwner = `

resource "tencentcloud_wedata_ops_task_owner" "wedata_ops_task_owner" {
}
`

package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsTaskOwnerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataOpsTaskOwner,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner", "owner_uin", "100029411056;100042282926"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner", "project_id", "2430455587205529600"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner", "task_id", "20251009144419600"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataOpsTaskOwnerUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner", "owner_uin", "100029411056"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner", "project_id", "2430455587205529600"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_task_owner.wedata_ops_task_owner", "task_id", "20251009144419600"),
				),
			},
		},
	})
}

const testAccWedataOpsTaskOwner = `

resource "tencentcloud_wedata_ops_task_owner" "wedata_ops_task_owner" {
    owner_uin  = "100029411056;100042282926"
    project_id = "2430455587205529600"
    task_id    = "20251009144419600"
}
`

const testAccWedataOpsTaskOwnerUp = `

resource "tencentcloud_wedata_ops_task_owner" "wedata_ops_task_owner" {
    owner_uin  = "100029411056"
    project_id = "2430455587205529600"
    task_id    = "20251009144419600"
}
`

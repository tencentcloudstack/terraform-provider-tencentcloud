package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCatTaskOpsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCatTaskOps,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cat_task_ops.task_ops", "id")),
			},
			{
				ResourceName:      "tencentcloud_cat_task_ops.task_ops",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCatTaskOps = `

resource "tencentcloud_cat_task_ops" "task_ops" {
  task_ids = 
}

`

package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsCompareTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsCompareTask,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_compare_task.compare_task", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_compare_task.compare_task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsCompareTask = `

resource "tencentcloud_dts_compare_task" "compare_task" {
  job_id = &lt;nil&gt;
  task_name = &lt;nil&gt;
  object_mode = &lt;nil&gt;
  objects {
		object_mode = &lt;nil&gt;
		object_items {
			db_name = &lt;nil&gt;
			db_mode = &lt;nil&gt;
			schema_name = &lt;nil&gt;
			table_mode = &lt;nil&gt;
			tables {
				table_name = &lt;nil&gt;
			}
			view_mode = &lt;nil&gt;
			views {
				view_name = &lt;nil&gt;
			}
		}

  }
  }

`

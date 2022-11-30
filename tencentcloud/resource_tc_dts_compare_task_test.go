package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDtsCompareTask_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsCompareTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_compare_task.compare_task", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_dts_compare_task.compareTask",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsCompareTask = `

resource "tencentcloud_dts_compare_task" "compare_task" {
  job_id = ""
  task_name = ""
  object_mode = ""
  objects {
			object_mode = ""
		object_items {
				db_name = ""
				db_mode = ""
				schema_name = ""
				table_mode = ""
			tables {
					table_name = ""
			}
				view_mode = ""
			views {
					view_name = ""
			}
		}

  }
  }

`

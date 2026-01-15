package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataDataBackfillPlanOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDataBackfillPlanOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_data_backfill_plan.example", "id"),
				),
			},
		},
	})
}

const testAccWedataDataBackfillPlanOperation = `
resource "tencentcloud_wedata_data_backfill_plan" "example" {
  project_id = "2430455587205529600"
  task_ids = [
    "20250625105147756"
  ]

  data_backfill_range_list {
    start_date           = "2026-01-14"
    end_date             = "2026-01-14"
    execution_start_time = "00:01"
    execution_end_time   = "23:59"
  }

  time_zone                         = "UTC+8"
  data_backfill_plan_name           = "tf-example"
  check_parent_type                 = "NONE"
  skip_event_listening              = true
  redefine_self_workflow_dependency = "no"
  redefine_parallel_num             = 2
  data_time_order                   = "NORMAL"
}
`

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
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_data_backfill_plan_operation.example", "id"),
				),
			},
		},
	})
}

const testAccWedataDataBackfillPlanOperation = `
resource "tencentcloud_wedata_data_backfill_plan_operation" "example" {
  project_id = "20241107221758402"
  task_ids = [
    "20250827115253729"
  ]

  data_backfill_range_list {
    start_date = "2025-09-02"
    end_date   = "2025-09-02"
  }

  skip_event_listening = true
}
`

package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCatProbeDataDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCatProbeData,
				Check:  resource.ComposeTestCheckFunc(
				//testAccCheckTencentCloudDataSourceID("data.tencentcloud_cat_probe_data.probe_data"),
				),
			},
		},
	})
}

const testAccDataSourceCatProbeData = `

data "tencentcloud_cat_probe_data" "probe_data" {
  begin_time = 1667923200000
  end_time = 1667996208428
  task_type = "AnalyzeTaskType_Network"
  sort_field = "ProbeTime"
  ascending = true
  selected_fields = ["terraform"]
  offset = 0
  limit = 20
  task_id = ["task-knare1mk"]
}

`

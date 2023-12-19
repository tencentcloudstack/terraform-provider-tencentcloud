package cat_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCatProbeDataDataSource -v
func TestAccTencentCloudCatProbeDataDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCatProbeData,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cat_probe_data.probe_data"),
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

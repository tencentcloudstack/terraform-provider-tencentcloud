package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCatProbedataDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCatProbedata,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cat_probedata.probe_data"),
				),
			},
		},
	})
}

const testAccDataSourceCatProbedata = `

data "tencentcloud_cat_probe_data" "probe_data" {
  begin_time = ""
  end_time = ""
  task_type = ""
  sort_field = "ProbeTime"
  ascending = ""
  selected_fields = ""
  offset = ""
  limit = ""
  task_i_d = ""
  operators = ""
  districts = ""
  error_types = ""
  city = ""
  code = ""
  }

`

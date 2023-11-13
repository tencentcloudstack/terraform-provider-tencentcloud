package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainTopSpaceTableTimeSeriesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainTopSpaceTableTimeSeriesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_top_space_table_time_series.top_space_table_time_series")),
			},
		},
	})
}

const testAccDbbrainTopSpaceTableTimeSeriesDataSource = `

data "tencentcloud_dbbrain_top_space_table_time_series" "top_space_table_time_series" {
  instance_id = ""
  limit = 
  sort_by = ""
  start_date = ""
  end_date = ""
  product = ""
  }

`

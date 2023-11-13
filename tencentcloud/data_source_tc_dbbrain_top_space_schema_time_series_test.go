package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainTopSpaceSchemaTimeSeriesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainTopSpaceSchemaTimeSeriesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_top_space_schema_time_series.top_space_schema_time_series")),
			},
		},
	})
}

const testAccDbbrainTopSpaceSchemaTimeSeriesDataSource = `

data "tencentcloud_dbbrain_top_space_schema_time_series" "top_space_schema_time_series" {
  instance_id = ""
  limit = 
  sort_by = ""
  start_date = ""
  end_date = ""
  product = ""
  }

`

package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDbbrainTopSpaceTableTimeSeriesObject = "data.tencentcloud_dbbrain_top_space_table_time_series.top_space_table_time_series"

func TestAccTencentCloudDbbrainTopSpaceTableTimeSeriesDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -2).In(loc).Format("2006-01-02")
	endTime := time.Now().In(loc).Format("2006-01-02")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainTopSpaceTableTimeSeriesDataSource, defaultDbBrainInstanceId, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testAccDbbrainTopSpaceTableTimeSeriesObject),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceTableTimeSeriesObject, "instance_id", defaultDbBrainInstanceId),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceTableTimeSeriesObject, "sort_by", "DataLength"),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceTableTimeSeriesObject, "product", "mysql"),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceTableTimeSeriesObject, "start_date", startTime),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceTableTimeSeriesObject, "end_date", endTime),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTableTimeSeriesObject, "top_space_table_time_series.#"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTableTimeSeriesObject, "top_space_table_time_series.0.table_name"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTableTimeSeriesObject, "top_space_table_time_series.0.table_schema"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTableTimeSeriesObject, "top_space_table_time_series.0.engine"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTableTimeSeriesObject, "top_space_table_time_series.0.series_data.0.series.#"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTableTimeSeriesObject, "top_space_table_time_series.0.series_data.0.series.0.metric"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTableTimeSeriesObject, "top_space_table_time_series.0.series_data.0.series.0.unit"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTableTimeSeriesObject, "top_space_table_time_series.0.series_data.0.series.0.values.#"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceTableTimeSeriesObject, "top_space_table_time_series.0.series_data.0.timestamp.#"),
				),
			},
		},
	})
}

const testAccDbbrainTopSpaceTableTimeSeriesDataSource = `

data "tencentcloud_dbbrain_top_space_table_time_series" "top_space_table_time_series" {
  instance_id = "%s"
  sort_by = "DataLength"
  start_date = "%s"
  end_date = "%s"
  product = "mysql"
}

`

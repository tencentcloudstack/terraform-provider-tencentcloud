package dbbrain_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDbbrainTopSpaceSchemaTimeSeriesObject = "data.tencentcloud_dbbrain_top_space_schema_time_series.top_space_schema_time_series"

func TestAccTencentCloudDbbrainTopSpaceSchemaTimeSeriesDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -2).In(loc).Format("2006-01-02")
	endTime := time.Now().In(loc).Format("2006-01-02")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainTopSpaceSchemaTimeSeriesDataSource, tcacctest.DefaultDbBrainInstanceId, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(testAccDbbrainTopSpaceSchemaTimeSeriesObject),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceSchemaTimeSeriesObject, "instance_id", tcacctest.DefaultDbBrainInstanceId),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceSchemaTimeSeriesObject, "sort_by", "DataLength"),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceSchemaTimeSeriesObject, "product", "mysql"),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceSchemaTimeSeriesObject, "start_date", startTime),
					resource.TestCheckResourceAttr(testAccDbbrainTopSpaceSchemaTimeSeriesObject, "end_date", endTime),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaTimeSeriesObject, "top_space_schema_time_series.#"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaTimeSeriesObject, "top_space_schema_time_series.0.table_schema"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaTimeSeriesObject, "top_space_schema_time_series.0.series_data.0.series.#"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaTimeSeriesObject, "top_space_schema_time_series.0.series_data.0.series.0.metric"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaTimeSeriesObject, "top_space_schema_time_series.0.series_data.0.series.0.unit"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaTimeSeriesObject, "top_space_schema_time_series.0.series_data.0.series.0.values.#"),
					resource.TestCheckResourceAttrSet(testAccDbbrainTopSpaceSchemaTimeSeriesObject, "top_space_schema_time_series.0.series_data.0.timestamp.#"),
				),
			},
		},
	})
}

const testAccDbbrainTopSpaceSchemaTimeSeriesDataSource = `

data "tencentcloud_dbbrain_top_space_schema_time_series" "top_space_schema_time_series" {
  instance_id = "%s"
  sort_by = "DataLength"
  start_date = "%s"
  end_date = "%s"
  product = "mysql"
  }

`

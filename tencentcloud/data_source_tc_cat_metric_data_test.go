package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCatMetricDataDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCatMetricDataDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cat_metric_data.metric_data")),
			},
		},
	})
}

const testAccCatMetricDataDataSource = `

data "tencentcloud_cat_metric_data" "metric_data" {
  analyze_task_type = ""
  metric_type = ""
  field = ""
  filter = ""
  group_by = ""
  filters = 
  }

`

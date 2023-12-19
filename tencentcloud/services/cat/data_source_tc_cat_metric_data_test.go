package cat_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCatMetricDataDataSource_basic -v
func TestAccTencentCloudCatMetricDataDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCatMetricDataDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cat_metric_data.metric_data"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cat_metric_data.metric_data", "metric_set"),
				),
			},
		},
	})
}

const testAccCatMetricDataDataSource = `
resource "tencentcloud_cat_task_set" "task_set" {
  interval = 1
  nodes = [
    "12136",
    "12137",
  ]
  parameters = jsonencode(
    {
      blackList         = ""
      filterIp          = 0
      grabBag           = 0
      ipType            = 0
      netDigOn          = 1
      netDnsNs          = ""
      netDnsOn          = 1
      netDnsQuerymethod = 1
      netDnsServer      = 2
      netDnsTimeout     = 5
      netIcmpActivex    = 0
      netIcmpActivexStr = ""
      netIcmpDataCut    = 1
      netIcmpInterval   = 0.5
      netIcmpNum        = 20
      netIcmpOn         = 1
      netIcmpSize       = 32
      netIcmpTimeout    = 20
      netTracertNum     = 30
      netTracertOn      = 1
      netTracertTimeout = 60
      whiteList         = ""
    }
  )
  tags          = {}
  task_category = 1
  task_type     = 5

  batch_tasks {
    name           = "terraform-test"
    target_address = "www.baidu.com"
  }
}

data "tencentcloud_cat_metric_data" "metric_data" {
  analyze_task_type = "AnalyzeTaskType_Network"
  metric_type = "gauge"
  field = "avg(\"ping_time\")"
  filters = [
    "\"host\" = 'www.baidu.com'",
    "time >= now()-1h",
  ]
  depends_on = [ tencentcloud_cat_task_set.task_set ]
}

`

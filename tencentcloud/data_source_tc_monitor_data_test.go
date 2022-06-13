package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonitorData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMonitorData(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_data.cvm_monitor_data"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_data.cvm_monitor_data",
						"list.#"),
				),
			},
		},
	})
}

func testAccDataSourceMonitorData() string {
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().Add(-time.Hour).In(loc).Format("2006-01-02T15:04:05+08:00")
	endTime := time.Now().Add(-40 * time.Minute).In(loc).Format("2006-01-02T15:04:05+08:00")

	return fmt.Sprintf(`
data "tencentcloud_instances" "instances" {
}

data "tencentcloud_monitor_data" "cvm_monitor_data" {
  namespace   = "QCE/CVM"
  metric_name = "CPUUsage"
  dimensions {
    name  = "InstanceId"
    value = data.tencentcloud_instances.instances.instance_list[0].instance_id
  }
  period     = 300
  start_time = "%s"
  end_time   = "%s"
}`, startTime, endTime)
}

package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonitorProductEvent(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMonitorProductEvent(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_product_event.cvm_event_data"),
				),
			},
		},
	})
}

func testAccDataSourceMonitorProductEvent() string {
	return fmt.Sprintf(`
data "tencentcloud_monitor_product_event" "cvm_event_data" {
  start_time      = %d
  is_alarm_config = 0
  product_name    = ["cvm"]
}`, time.Now().Unix()-86400)
}

package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudEbBusDataSource_basic -v
func TestAccTencentCloudEbBusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbBusDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_eb_bus.bus"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_bus.bus", "event_buses.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_bus.bus", "event_buses.0.add_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_bus.bus", "event_buses.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_bus.bus", "event_buses.0.event_bus_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_bus.bus", "event_buses.0.event_bus_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_bus.bus", "event_buses.0.mod_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_bus.bus", "event_buses.0.pay_mode"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_bus.bus", "event_buses.0.type"),
				),
			},
		},
	})
}

const testAccEbBusDataSource = `

resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}
data "tencentcloud_eb_bus" "bus" {
  order_by = "AddTime"
  order = "DESC"
  filters {
	values = ["Custom"]
	name = "Type"
  }
  depends_on = [ tencentcloud_eb_event_bus.foo ]
}

`

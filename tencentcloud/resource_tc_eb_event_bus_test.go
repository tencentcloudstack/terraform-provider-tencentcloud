package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudEbEventBusResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbEventBus,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_eb_event_bus.event_bus", "id")),
			},
			{
				ResourceName:      "tencentcloud_eb_event_bus.event_bus",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccEbEventBus = `

resource "tencentcloud_eb_event_bus" "event_bus" {
  event_bus_name = ""
  description = ""
  save_days = 
  enable_store = 
  tags = {
    "createdBy" = "terraform"
  }
}

`

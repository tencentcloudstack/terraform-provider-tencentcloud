package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsEventResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsEvent,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_event.event", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_event.event",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsEvent = `

resource "tencentcloud_mps_event" "event" {
  event_name = "you-event-name"
  description = "event description"
}

`

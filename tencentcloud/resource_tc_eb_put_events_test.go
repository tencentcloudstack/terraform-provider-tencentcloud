package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudEbPutEventsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbPutEvents,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_eb_put_events.put_events", "id")),
			},
			{
				ResourceName:      "tencentcloud_eb_put_events.put_events",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccEbPutEvents = `

resource "tencentcloud_eb_put_events" "put_events" {
  event_list {
		source = ""
		data = ""
		type = ""
		subject = ""
		time = 

  }
  event_bus_id = ""
}

`

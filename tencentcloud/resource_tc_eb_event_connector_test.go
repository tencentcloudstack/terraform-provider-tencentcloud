package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudEbEventConnectorResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbEventConnector,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_eb_event_connector.event_connector", "id")),
			},
			{
				ResourceName:      "tencentcloud_eb_event_connector.event_connector",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccEbEventConnector = `

resource "tencentcloud_eb_event_connector" "event_connector" {
  connection_description {
		resource_description = ""
		a_p_i_g_w_params {
			protocol = ""
			method = ""
		}
		ckafka_params {
			offset = ""
			topic_name = ""
		}
		d_t_s_params = 

  }
  event_bus_id = ""
  connection_name = ""
  description = ""
  enable = 
  type = ""
}

`

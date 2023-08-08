package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudEbEventTransformResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbEventTransform,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_eb_event_transform.eb_transform", "id")),
			},
			{
				ResourceName:      "tencentcloud_eb_event_transform.eb_transform",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccEbEventTransform = `

resource "tencentcloud_eb_event_transform" "eb_transform" {
  event_bus_id = ""
  rule_id = ""
  transformations {
		extraction {
			extraction_input_path = ""
			format = ""
			text_params {
				separator = ""
				regex = ""
			}
		}
		etl_filter {
			filter = ""
		}
		transform {
			output_structs {
				key = ""
				value = ""
				value_type = ""
			}
		}

  }
}

`

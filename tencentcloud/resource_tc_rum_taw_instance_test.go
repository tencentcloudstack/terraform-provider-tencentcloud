package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudRumTawInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumTawInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_rum_taw_instance.taw_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_rum_taw_instance.tawInstance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRumTawInstance = `

resource "tencentcloud_rum_taw_instance" "taw_instance" {
  area_id = ""
  charge_type = ""
  data_retention_days = ""
  instance_name = ""
  tags {
			key = ""
			value = ""

  }
  instance_desc = ""
  count_num = ""
  period_retain = ""
  buying_channel = ""
            }

`

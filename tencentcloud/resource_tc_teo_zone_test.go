package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTeoZone_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoZone,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_zone.zone", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_zone.zone",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoZone = `

resource "tencentcloud_teo_zone" "zone" {
  name           = ""
  plan_type      = ""
  type           = ""
  paused         = ""
  vanity_name_servers {
    switch  = ""
    servers = ""

  }
  cname_speed_up = ""
  tags           = {
    "createdBy" = "terraform"
  }
}

`

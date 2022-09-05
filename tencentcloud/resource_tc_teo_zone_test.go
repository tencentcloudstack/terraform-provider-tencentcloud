package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTeoZone_basic(t *testing.T) {
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
  name           = "sfurnace.work"
  plan_type      = "ent_cm_with_bot"
  type           = "full"
  paused         = false
  cname_speed_up = "enabled"

  #  vanity_name_servers {
  #    switch  = "on"
  #    servers = ["2.2.2.2"]
  #  }
}

`

package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoZoneResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoZone,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_zone.zone", "id")),
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
    zone_name = &lt;nil&gt;
  plan_type = &lt;nil&gt;
          type = &lt;nil&gt;
  paused = &lt;nil&gt;
      vanity_name_servers {
		switch = &lt;nil&gt;
		servers = &lt;nil&gt;

  }
    cname_speed_up = &lt;nil&gt;
    tags {
		tag_key = &lt;nil&gt;
		tag_value = &lt;nil&gt;

  }
      tags = {
    "createdBy" = "terraform"
  }
}

`

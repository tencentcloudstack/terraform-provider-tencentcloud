package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTeoLoadBalancing_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoLoadBalancing,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.originGroup", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_origin_group.originGroup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoLoadBalancing = `

resource "tencentcloud_teo_load_balancing" "load_balancing" {
#  backup_origin_group_id = "origin-a499ca4b-3721-11ed-b9c1-5254005a52aa"
  host                   = "www.toutiao2.com"
  origin_group_id        = "origin-4f8a30b2-3720-11ed-b66b-525400dceb86"
  status                 = "online"
  tags                   = {}
  ttl                    = 600
  type                   = "proxied"
  zone_id                = "zone-297z8rf93cfw"
}

`

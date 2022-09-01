package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTeoOriginGroup_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoOriginGroup,
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

const testAccTeoOriginGroup = `

locals {
  group0 = [
    {
      "record" = "1.1.1.1"
      "port"   = 80
      "weight" = 30
    }, {
      "record" = "2.2.2.2"
      "port"   = 443
      "weight" = 70
    }
  ]
}

resource "tencentcloud_teo_origin_group" "origin_group" {
  zone_id     = tencentcloud_teo_zone.zone.id
  origin_name = "group0"
  origin_type = "self"
  type        = "weight"

  dynamic "record" {
    for_each = local.group0
    content {
      record = record.value["record"]
      port   = record.value["port"]
      weight = record.value["weight"]
      area   = []
    }
  }
}

`

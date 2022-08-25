package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTeoOriginGroup_basic(t *testing.T) {
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

resource "tencentcloud_teo_origin_group" "originGroup" {
  origin_name = "test"
  type        = "weight"
  record {
    record  = "20160527-10003318.cos.ap-shanghai.myqcloud.com"
    area    = []
    weight  = 100
    port    = 0
    private = false

  }
  zone_id     = "zone-27mypfc1vr7d"
  origin_type = "cos"
}

`

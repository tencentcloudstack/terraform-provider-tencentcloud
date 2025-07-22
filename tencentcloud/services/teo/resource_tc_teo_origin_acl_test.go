package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoOriginAclResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoOriginAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_acl.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_acl.example", "zone_id"),
				),
			},
			{
				Config: testAccTeoOriginAclUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_acl.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_acl.example", "zone_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_origin_acl.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoOriginAcl = `
resource "tencentcloud_teo_origin_acl" "example" {
  zone_id        = "zone-3edjdliiw3he"
  l7_hosts = [
    "1.makn.cn",
    "2.makn.cn",
    "3.makn.cn",
    "4.makn.cn",
    "5.makn.cn",
    "6.makn.cn",
  ]

  l4_proxy_ids = [
    "sid-3edjfy5n10wh",
    "sid-3edjg2kml7pg",
    "sid-3edjg6t8dw78",
    "sid-3edjgc30nbgx",
  ]

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
`

const testAccTeoOriginAclUpdate = `
resource "tencentcloud_teo_origin_acl" "example" {
  zone_id        = "zone-3edjdliiw3he"
  l7_hosts = [
    "1.makn.cn",
    "2.makn.cn",
    "3.makn.cn",
  ]

  l4_proxy_ids = [
    "sid-3edjg6t8dw78",
    "sid-3edjgc30nbgx",
  ]

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
`

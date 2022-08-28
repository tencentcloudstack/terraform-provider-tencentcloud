package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTeoDnsSec_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsSec,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_dns_sec.dnsSec", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_dns_sec.dnsSec",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoDnsSec = `

resource "tencentcloud_teo_dns_sec" "dnsSec" {
  zone_id = ""
  status  = ""
}
`

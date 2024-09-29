package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -test.run TestAccTencentCloudTeoSecurityIpGroupResource_basic -v
func TestAccTencentCloudTeoSecurityIpGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoSecurityIpGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_ip_group.teo_security_ip_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_ip_group.teo_security_ip_group", "ip_group.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_ip_group.teo_security_ip_group", "ip_group.0.name", "aaaaa"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_ip_group.teo_security_ip_group", "ip_group.0.content.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_security_ip_group.teo_security_ip_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoSecurityIpGroupUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_ip_group.teo_security_ip_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_ip_group.teo_security_ip_group", "ip_group.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_ip_group.teo_security_ip_group", "ip_group.0.name", "bbbbb"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_ip_group.teo_security_ip_group", "ip_group.0.content.#", "3"),
				),
			},
		},
	})
}

const testAccTeoSecurityIpGroup = `

resource "tencentcloud_teo_security_ip_group" "teo_security_ip_group" {
    zone_id = "zone-2qtuhspy7cr6"

    ip_group {
        content  = [
            "10.1.1.1",
            "10.1.1.2",
        ]
        name     = "aaaaa"
    }
}
`

const testAccTeoSecurityIpGroupUp = `

resource "tencentcloud_teo_security_ip_group" "teo_security_ip_group" {
    zone_id = "zone-2qtuhspy7cr6"

    ip_group {
        content  = [
            "10.1.1.1",
            "10.1.1.2",
            "10.1.1.3",
        ]
        name     = "bbbbb"
    }
}
`

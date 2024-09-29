package teo

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoSecurityIpGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoSecurityIpGroup,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_security_ip_group.teo_security_ip_group", "id")),
		}, {
			ResourceName:      "tencentcloud_teo_security_ip_group.teo_security_ip_group",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccTeoSecurityIpGroup = `

resource "tencentcloud_teo_security_ip_group" "teo_security_ip_group" {
  ip_group = {
  }
}
`

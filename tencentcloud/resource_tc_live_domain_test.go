package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveDomain,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_domain.domain", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_domain.domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveDomain = `

resource "tencentcloud_live_domain" "domain" {
  domain_name = ""
  domain_type = 
  play_type = 
  is_delay_live = 
  is_mini_program_live = 
  verify_owner_type = ""
}

`

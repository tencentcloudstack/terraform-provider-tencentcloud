package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveDomainCertResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveDomainCert,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_domain_cert.domain_cert", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_domain_cert.domain_cert",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveDomainCert = `

resource "tencentcloud_live_domain_cert" "domain_cert" {
  domain_name = "5000.livepush.play.com"
  type = "Formal"
}

`

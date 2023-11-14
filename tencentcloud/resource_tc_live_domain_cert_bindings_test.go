package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveDomainCertBindingsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveDomainCertBindings,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_domain_cert_bindings.domain_cert_bindings", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_domain_cert_bindings.domain_cert_bindings",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveDomainCertBindings = `

resource "tencentcloud_live_domain_cert_bindings" "domain_cert_bindings" {
  domain_infos {
		domain_name = "abc.com"
		status = 1

  }
  cloud_cert_id = "123"
  certificate_public_key = "xxx"
  certificate_private_key = "xxx"
  certificate_alias = "adc"
}

`

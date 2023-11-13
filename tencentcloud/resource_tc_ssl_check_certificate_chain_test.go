package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslCheckCertificateChainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslCheckCertificateChain,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_check_certificate_chain.check_certificate_chain", "id")),
			},
			{
				ResourceName:      "tencentcloud_ssl_check_certificate_chain.check_certificate_chain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSslCheckCertificateChain = `

resource "tencentcloud_ssl_check_certificate_chain" "check_certificate_chain" {
  certificate_chain = ""
}

`

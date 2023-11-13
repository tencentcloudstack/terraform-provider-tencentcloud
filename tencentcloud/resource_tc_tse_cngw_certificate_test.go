package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseCngwCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwCertificate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_certificate.cngw_certificate", "id")),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_certificate.cngw_certificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTseCngwCertificate = `

resource "tencentcloud_tse_cngw_certificate" "cngw_certificate" {
  gateway_id = ""
  bind_domains = 
  cert_id = ""
  name = ""
  key = ""
  crt = ""
}

`

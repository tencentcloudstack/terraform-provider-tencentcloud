package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslReplaceCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslReplaceCertificate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_replace_certificate.replace_certificate", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_replace_certificate.replace_certificate", "certificate_id", "8hUkH3xC"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_replace_certificate.replace_certificate", "valid_type", "DNS_AUTO"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_replace_certificate.replace_certificate", "csr_type", "online"),
				),
			},
		},
	})
}

const testAccSslReplaceCertificate = `

resource "tencentcloud_ssl_replace_certificate" "replace_certificate" {
  certificate_id = "8hUkH3xC"
  valid_type = "DNS_AUTO"
  csr_type = "online"
}

`

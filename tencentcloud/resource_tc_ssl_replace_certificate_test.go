package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_replace_certificate.replace_certificate", "id")),
			},
			{
				ResourceName:      "tencentcloud_ssl_replace_certificate.replace_certificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSslReplaceCertificate = `

resource "tencentcloud_ssl_replace_certificate" "replace_certificate" {
  certificate_id = ""
  valid_type = ""
  csr_type = ""
  csr_content = ""
  csrkey_password = ""
  reason = ""
  cert_c_s_r_encrypt_algo = ""
  cert_c_s_r_key_parameter = ""
}

`

package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslCommitCertificateInformationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslCommitCertificateInformation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_commit_certificate_information.commit_certificate_information", "id")),
			},
		},
	})
}

const testAccSslCommitCertificateInformation = `

resource "tencentcloud_ssl_pay_certificate" "ssl" {
  product_id = 33
  domain_num = 1
  alias      = "test-ssl-update"
  project_id = 0
  wait_commit_flag = true
  information {
    csr_type              = "online"
    certificate_domain    = "www.domain.com"
    organization_name     = "test-update"
    organization_division = "test"
    organization_address  = "test"
    organization_country  = "CN"
    organization_city     = "test"
    organization_region   = "test"
    postal_code           = "0755"
    phone_area_code       = "0755"
    phone_number          = "12345678901"
    verify_type           = "DNS"
    admin_first_name      = "test"
    admin_last_name       = "test"
    admin_phone_num       = "12345678901"
    admin_email           = "test@tencent.com"
    admin_position        = "dev"
    contact_first_name    = "test"
    contact_last_name     = "test"
    contact_email         = "test@tencent.com"
    contact_number        = "12345678901"
    contact_position      = "dev"
  }
}
resource "tencentcloud_ssl_commit_certificate_information" "commit_certificate_information" {
  product_id = 33
  certificate_id           = tencentcloud_ssl_pay_certificate.ssl.certificate_id
}

`

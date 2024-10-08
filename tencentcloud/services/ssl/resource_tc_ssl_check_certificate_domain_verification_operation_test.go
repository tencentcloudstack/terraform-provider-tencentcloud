package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSslCheckCertificateDomainVerificationOperationResource_basic -v
func TestAccTencentCloudSslCheckCertificateDomainVerificationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslCheckCertificateDomainVerification,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_check_certificate_domain_verification_operation.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_check_certificate_domain_verification_operation.example", "certificate_id"),
				),
			},
		},
	})
}

const testAccSslCheckCertificateDomainVerification = `
resource "tencentcloud_ssl_check_certificate_domain_verification_operation" "example" {
  certificate_id = "6BE701Jx"
}
`

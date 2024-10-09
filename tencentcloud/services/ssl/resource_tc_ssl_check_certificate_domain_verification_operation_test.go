package ssl_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixSslCheckCertificateDomainVerificationOperationResource_basic -v
func TestAccTencentCloudNeedFixSslCheckCertificateDomainVerificationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslCheckCertificateDomainVerificationOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_check_certificate_domain_verification_operation.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_check_certificate_domain_verification_operation.example", "certificate_id"),
				),
			},
		},
	})
}

const testAccSslCheckCertificateDomainVerificationOperation = `
resource "tencentcloud_ssl_check_certificate_domain_verification_operation" "example" {
  certificate_id = "6BE701Jx"
}

`

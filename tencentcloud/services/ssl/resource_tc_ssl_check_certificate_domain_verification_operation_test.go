package ssl

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudSslCheckCertificateDomainVerificationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccSslCheckCertificateDomainVerificationOperation,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_check_certificate_domain_verification_operation.ssl_check_certificate_domain_verification_operation", "id")),
		}, {
			ResourceName:      "tencentcloud_ssl_check_certificate_domain_verification_operation.ssl_check_certificate_domain_verification_operation",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccSslCheckCertificateDomainVerificationOperation = `

resource "tencentcloud_ssl_check_certificate_domain_verification_operation" "ssl_check_certificate_domain_verification_operation" {
}
`

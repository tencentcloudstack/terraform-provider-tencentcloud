package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslCompleteCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslCompleteCertificate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_complete_certificate_operation.complete_certificate", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_complete_certificate_operation.complete_certificate", "certificate_id", "709ahm2q"),
				),
			},
		},
	})
}

const testAccSslCompleteCertificate = `

resource "tencentcloud_ssl_complete_certificate_operation" "complete_certificate" {
  certificate_id = "709ahm2q"
}

`

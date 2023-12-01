package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslCompleteCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
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

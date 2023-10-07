package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslCompleteCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslCompleteCertificate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_complete_certificate.complete_certificate", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_complete_certificate.complete_certificate", "certificate_id", "9Bfe1IBR"),
				),
			},
			{
				ResourceName:      "tencentcloud_ssl_complete_certificate.complete_certificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSslCompleteCertificate = `

resource "tencentcloud_ssl_complete_certificate" "complete_certificate" {
  certificate_id = "9Bfe1IBR"
}

`

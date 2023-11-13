package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_complete_certificate.complete_certificate", "id")),
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
  certificate_id = ""
}

`

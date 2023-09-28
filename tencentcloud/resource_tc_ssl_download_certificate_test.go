package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslDownloadCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDownloadCertificate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_download_certificate.download_certificate", "id")),
			},
		},
	})
}

const testAccSslDownloadCertificate = `

resource "tencentcloud_ssl_download_certificate" "download_certificate" {
  certificate_id = "8x1eUSSl"
  output_path = "./"
}

`

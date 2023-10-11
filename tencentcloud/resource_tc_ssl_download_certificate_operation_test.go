package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDownloadCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDownloadCertificate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_download_certificate_operation.download_certificate", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_download_certificate_operation.download_certificate", "certificate_id", "8x1eUSSl"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_download_certificate_operation.download_certificate", "output_path"),
				),
			},
		},
	})
}

const testAccSslDownloadCertificate = `

resource "tencentcloud_ssl_download_certificate_operation" "download_certificate" {
  certificate_id = "8x1eUSSl"
  output_path = "./"
}

`

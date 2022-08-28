package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTeoHostCertificate_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoHostCertificate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_host_certificate.hostCertificate", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_host_certificate.hostCertificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoHostCertificate = `

resource "tencentcloud_teo_host_certificate" "hostCertificate" {
  zone_id = ""
  host    = ""
  cert_info {
    cert_id = ""
    status  = ""
  }
}

`

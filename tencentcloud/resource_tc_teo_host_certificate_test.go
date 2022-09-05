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
					resource.TestCheckResourceAttrSet("tencentcloud_teo_host_certificate.host_certificate", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_host_certificate.host_certificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoHostCertificate = `

resource "tencentcloud_teo_host_certificate" "host_certificate" {
  zone_id = tencentcloud_teo_zone.zone.id
  host    = tencentcloud_teo_dns_record.dns_record.name

  cert_info {
    cert_id = "yqWPPbs7"
    status  = "deployed"
  }
}

`

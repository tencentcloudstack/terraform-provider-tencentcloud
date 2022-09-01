package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTeoDefaultCertificate_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDefaultCertificate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_default_certificate.default_certificate", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_default_certificate.defaultCertificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoDefaultCertificate = `

resource "tencentcloud_teo_default_certificate" "default_certificate" {
  zone_id = tencentcloud_teo_zone.zone.id

  cert_info {
    cert_id = "teo-28i46c1gtmkl"
    status  = "deployed"
  }
}

`

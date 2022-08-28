package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTeoDefaultCertificate_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDefaultCertificate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_default_certificate.defaultCertificate", "id"),
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

resource "tencentcloud_teo_default_certificate" "defaultCertificate" {
  zone_id = ""
  cert_id = ""
}

`

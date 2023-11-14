package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoDefaultCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDefaultCertificate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_default_certificate.default_certificate", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_default_certificate.default_certificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoDefaultCertificate = `

resource "tencentcloud_teo_default_certificate" "default_certificate" {
  zone_id = &lt;nil&gt;
  cert_info {
		cert_id = &lt;nil&gt;
		status = &lt;nil&gt;

  }
}

`

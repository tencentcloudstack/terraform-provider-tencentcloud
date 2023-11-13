package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoHostCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoHostCertificate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_host_certificate.host_certificate", "id")),
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
  zone_id = &lt;nil&gt;
  host = &lt;nil&gt;
  cert_info {
		cert_id = &lt;nil&gt;
		status = &lt;nil&gt;

  }
}

`

package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoCertificate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate.certificate", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_certificate.certificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoCertificate = `

resource "tencentcloud_teo_certificate" "certificate" {
  zone_id = ""
  hosts = 
  server_cert_info {
		cert_id = ""
		alias = ""
		type = ""
		expire_time = ""
		deploy_time = ""
		sign_algo = ""
		common_name = ""

  }
  mode = ""
}

`

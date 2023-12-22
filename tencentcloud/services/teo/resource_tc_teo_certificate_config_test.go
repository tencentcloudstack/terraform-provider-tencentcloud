package teo_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTeoCertificateConfigResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoCertificateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate_config.certificate", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_certificate_config.certificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoCertificateConfig = testAccTeoZone + `

resource "tencentcloud_teo_certificate_config" "certificate" {
    host    = ` + "test." + `var.zone_name
    mode    = "disable"
    zone_id = tencentcloud_teo_zone.basic.id

    server_cert_info {
        alias       = "EdgeOne default"
        cert_id     = "teo-2o1tfutpnb6l"
        common_name = var.zone_name
        deploy_time = "2023-09-27T11:54:47Z"
        expire_time = "2023-12-26T06:38:47Z"
        sign_algo   = "RSA 2048"
        type        = "default"
    }
}

`

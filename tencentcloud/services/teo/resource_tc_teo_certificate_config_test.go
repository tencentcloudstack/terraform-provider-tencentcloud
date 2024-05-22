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
					resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate_config.certificate", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate_config.certificate", "host"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "mode", "sslcert"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.alias", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.cert_id", "EEIqXrZt"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.common_name"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.sign_algo", "RSA 2048"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.type", "managed"),
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

resource "tencentcloud_teo_ownership_verify" "ownership_verify" {
  domain = var.zone_name

  depends_on = [ tencentcloud_teo_zone.basic ]
}

resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
  zone_id     = tencentcloud_teo_zone.basic.id
  domain_name = "test.tf-teo.xyz"

  origin_info {
    origin      = "150.109.8.1"
    origin_type = "IP_DOMAIN"
  }

  depends_on = [ tencentcloud_teo_ownership_verify.ownership_verify ]
}

resource "tencentcloud_teo_certificate_config" "certificate" {
  host    = format("test.%s", var.zone_name)
  mode    = "sslcert"
  zone_id = tencentcloud_teo_zone.basic.id

  server_cert_info {
    alias       = "terraform_test"
    cert_id     = "EEIqXrZt"
    common_name = var.zone_name
    //deploy_time = "2024-04-22T10:34:13Z"
    //expire_time = "2025-04-22T23:59:59Z"
    sign_algo   = "RSA 2048"
    type        = "managed"
  }

  depends_on = [tencentcloud_teo_acceleration_domain.acceleration_domain]
}

`

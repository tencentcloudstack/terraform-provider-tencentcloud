package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttDeviceCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMqttDeviceCertificate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "device_certificate"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "ca_sn"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "client_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "format"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "status"),
				),
			},
			{
				Config: testAccMqttDeviceCertificateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "device_certificate"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "ca_sn"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "client_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "format"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_device_certificate.example", "status"),
				),
			},
			{
				ResourceName:      "tencentcloud_mqtt_device_certificate.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMqttDeviceCertificate = `
resource "tencentcloud_mqtt_device_certificate" "example" {
  instance_id        = "mqtt-zxjwkr98"
  device_certificate = ""
  ca_sn              = ""
  client_id          = ""
  format             = ""
  status             = "ACTIVE"
}
`

const testAccMqttDeviceCertificateUpdate = `
resource "tencentcloud_mqtt_device_certificate" "example" {
  instance_id        = "mqtt-zxjwkr98"
  device_certificate = ""
  ca_sn              = ""
  client_id          = ""
  format             = ""
  status             = "ACTIVE"
}
`

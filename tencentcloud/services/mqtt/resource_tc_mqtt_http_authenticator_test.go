package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttHttpAuthenticatorResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMqttHttpAuthenticator,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "endpoint"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "concurrency"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "method"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "connect_timeout"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "read_timeout"),
				),
			},
			{
				Config: testAccMqttHttpAuthenticatorUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "endpoint"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "concurrency"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "method"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "connect_timeout"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.example", "read_timeout"),
				),
			},
			{
				ResourceName:      "tencentcloud_mqtt_http_authenticator.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMqttHttpAuthenticator = `
resource "tencentcloud_mqtt_http_authenticator" "example" {
  instance_id     = "mqtt-zxjwkr98"
  endpoint        = "https://example.com"
  concurrency     = 8
  method          = "POST"
  status          = "open"
  remark          = "Remark."
  connect_timeout = 10
  read_timeout    = 10
  header {
    key   = "Content-type"
    value = "application/json"
  }

  body {
    key   = "bodyKey"
    value = "bodyValue"
  }
}
`

const testAccMqttHttpAuthenticatorUpdate = `
resource "tencentcloud_mqtt_http_authenticator" "example" {
  instance_id     = "mqtt-zxjwkr98"
  endpoint        = "https://example.com"
  concurrency     = 8
  method          = "GET"
  status          = "close"
  remark          = "Remark update."
  connect_timeout = 10
  read_timeout    = 10
  header {
    key   = "Content-type"
    value = "application/json"
  }

  body {
    key   = "bodyKey"
    value = "bodyValue"
  }
}
`

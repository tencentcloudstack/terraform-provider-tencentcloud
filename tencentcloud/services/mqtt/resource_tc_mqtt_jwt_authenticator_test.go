package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttJwtAuthenticatorResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMqttJwtAuthenticator,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "algorithm"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "from"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "secret"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "remark"),
				),
			},
			{
				Config: testAccMqttJwtAuthenticatorUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "algorithm"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "from"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "secret"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwt_authenticator.example", "remark"),
				),
			},
			{
				ResourceName:      "tencentcloud_mqtt_jwt_authenticator.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMqttJwtAuthenticator = `
resource "tencentcloud_mqtt_jwt_authenticator" "example" {
  instance_id = "mqtt-zxjwkr98"
  algorithm   = "hmac-based"
  from        = "password"
  secret      = "your secret content"
  status      = "open"
  remark      = "Remark."
}
`

const testAccMqttJwtAuthenticatorUpdate = `
resource "tencentcloud_mqtt_jwt_authenticator" "example" {
  instance_id = "mqtt-zxjwkr98"
  algorithm   = "hmac-based"
  from        = "password"
  secret      = "your secret content"
  status      = "open"
  remark      = "Remark update."
}
`

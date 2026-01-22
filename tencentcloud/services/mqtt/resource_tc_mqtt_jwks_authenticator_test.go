package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttJwksAuthenticatorResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccMqttJwksAuthenticator,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mqtt_jwks_authenticator.mqtt_jwks_authenticator", "id")),
		}, {
			ResourceName:      "tencentcloud_mqtt_jwks_authenticator.mqtt_jwks_authenticator",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccMqttJwksAuthenticator = `

resource "tencentcloud_mqtt_jwks_authenticator" "mqtt_jwks_authenticator" {
}
`

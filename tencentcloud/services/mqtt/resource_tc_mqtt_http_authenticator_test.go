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
		Steps: []resource.TestStep{{
			Config: testAccMqttHttpAuthenticator,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mqtt_http_authenticator.mqtt_http_authenticator", "id")),
		}, {
			ResourceName:      "tencentcloud_mqtt_http_authenticator.mqtt_http_authenticator",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccMqttHttpAuthenticator = `

resource "tencentcloud_mqtt_http_authenticator" "mqtt_http_authenticator" {
  header = {
  }
  body = {
  }
}
`

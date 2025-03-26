package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttApplyRegistrationCodeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccMqttApplyRegistrationCodeDataSource,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.tencentcloud_mqtt_apply_registration_code.example", "id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_mqtt_apply_registration_code.example", "instance_id"),
			),
		}},
	})
}

const testAccMqttApplyRegistrationCodeDataSource = `
data "tencentcloud_mqtt_apply_registration_code" "example" {
  instance_id = "mqtt-zxjwkr98"
}
`

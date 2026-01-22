package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttTopicsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccMqttTopicsDataSource,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.tencentcloud_mqtt_topics.example", "id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_mqtt_topics.example", "instance_id"),
			),
		}},
	})
}

const testAccMqttTopicsDataSource = `
data "tencentcloud_mqtt_topics" "example" {
  instance_id = "mqtt-g4qgr3gx"
}
`

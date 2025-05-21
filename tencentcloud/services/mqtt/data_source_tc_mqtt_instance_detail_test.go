package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttInstanceDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccMqttInstanceDetailDataSource,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.tencentcloud_mqtt_instance_detail.example", "id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_mqtt_instance_detail.example", "instance_id"),
			),
		}},
	})
}

const testAccMqttInstanceDetailDataSource = `
data "tencentcloud_mqtt_instance_detail" "example" {
  instance_id = "mqtt-kngmpg9p"
}
`

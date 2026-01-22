package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccMqttInstancesDataSource,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.tencentcloud_mqtt_instances.example", "id"),
			),
		}},
	})
}

const testAccMqttInstancesDataSource = `
data "tencentcloud_mqtt_instances" "example" {}
`

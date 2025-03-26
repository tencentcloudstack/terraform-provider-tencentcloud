package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttUserResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMqttUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_user.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_user.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_user.example", "username"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_user.example", "password"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_user.example", "remark"),
				),
			},
			{
				Config: testAccMqttUserUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_user.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_user.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_user.example", "username"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_user.example", "password"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_user.example", "remark"),
				),
			},
			{
				ResourceName:      "tencentcloud_mqtt_user.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMqttUser = `
resource "tencentcloud_mqtt_user" "example" {
  instance_id = "mqtt-zxjwkr98"
  username    = "tf-example"
  password    = "Password@123"
  remark      = "Remark."
}
`

const testAccMqttUserUpdate = `
resource "tencentcloud_mqtt_user" "example" {
  instance_id = "mqtt-zxjwkr98"
  username    = "tf-example"
  password    = "Password@123"
  remark      = "Remark update."
}
`

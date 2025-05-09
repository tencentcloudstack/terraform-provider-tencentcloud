package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttAuthorizationPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMqttAuthorizationPolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "policy_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "policy_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "priority"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "effect"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "actions"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "retain"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "qos"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "resources"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "username"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "client_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "remark"),
				),
			},
			{
				Config: testAccMqttAuthorizationPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "policy_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "policy_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "priority"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "effect"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "actions"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "retain"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "qos"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "resources"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "username"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "client_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_authorization_policy.example", "remark"),
				),
			},
			{
				ResourceName:      "tencentcloud_mqtt_authorization_policy.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMqttAuthorizationPolicy = `
resource "tencentcloud_mqtt_authorization_policy" "example" {
  instance_id    = "mqtt-g4qgr3gx"
  policy_name    = "tf-example"
  policy_version = 1
  priority       = 10
  effect         = "allow"
  actions        = "connect,pub,sub"
  retain         = 3
  qos            = "0,1,2"
  resources      = "topic-demo"
  username       = "*root*"
  client_id      = "client"
  ip             = "192.168.1.1"
  remark         = "policy remark."
}
`

const testAccMqttAuthorizationPolicyUpdate = `
resource "tencentcloud_mqtt_authorization_policy" "example" {
  instance_id    = "mqtt-g4qgr3gx"
  policy_name    = "tf-example-update"
  policy_version = 1
  priority       = 11
  effect         = "deny"
  actions        = "pub,sub"
  retain         = 2
  qos            = "1,2"
  resources      = "topic-demo"
  username       = "root"
  client_id      = "*client*"
  ip             = "192.168.1.0/24"
  remark         = "policy remark update."
}
`

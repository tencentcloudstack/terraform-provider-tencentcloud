package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttMessageEnrichmentRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMqttMessageEnrichmentRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "rule_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "condition"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "actions"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "priority"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "rule_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "update_time"),
				),
			},
			{
				ResourceName:      "tencentcloud_mqtt_message_enrichment_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudMqttMessageEnrichmentRuleResource_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMqttMessageEnrichmentRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "rule_name"),
					resource.TestCheckResourceAttr("tencentcloud_mqtt_message_enrichment_rule.example", "priority", "10"),
					resource.TestCheckResourceAttr("tencentcloud_mqtt_message_enrichment_rule.example", "status", "1"),
				),
			},
			{
				Config: testAccMqttMessageEnrichmentRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.example", "rule_name"),
					resource.TestCheckResourceAttr("tencentcloud_mqtt_message_enrichment_rule.example", "priority", "5"),
					resource.TestCheckResourceAttr("tencentcloud_mqtt_message_enrichment_rule.example", "status", "1"),
				),
			},
		},
	})
}

const testAccMqttMessageEnrichmentRule = `
resource "tencentcloud_mqtt_message_enrichment_rule" "example" {
  instance_id = "mqtt-zxje8zdd"
  rule_name   = "tf-example"
  condition {
    username  = "user*"
    client_id = "clientDemo"
    topic     = "topicDemo"
  }

  actions {
    message_expiry_interval = 3600
    response_topic          = "topicDemo"
    correlation_data        = "correlationData"
    user_property {
      key   = "key"
      value = "value"
    }
  }
  priority = 10
  status   = 1
  remark   = "remark."
}
`

const testAccMqttMessageEnrichmentRuleUpdate = `
resource "tencentcloud_mqtt_message_enrichment_rule" "example" {
  instance_id = "mqtt-zxje8zdd"
  rule_name   = "tf-example"
  condition {
    username  = "user*"
    client_id = "clientDemo"
    topic     = "topicDemo"
  }

  actions {
    correlation_data        = "correlationData"
  }
  priority = 1
  status   = 2
  remark   = "remark update."
}
`

Provides a resource to create a MQTT message enrichment rule

Example Usage

```hcl
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
```

Import

MQTT message enrichment rule can be imported using the instanceId#ruleId, e.g.

```
terraform import tencentcloud_mqtt_message_enrichment_rule.example mqtt-zxje8zdd#34
```

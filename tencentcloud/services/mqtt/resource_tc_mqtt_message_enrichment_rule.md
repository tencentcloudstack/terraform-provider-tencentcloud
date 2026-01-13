Provides a resource to create a MQTT message enrichment rule.

Message enrichment rules allow you to dynamically add or modify message properties when messages are published, such as setting message expiration time, response topic, correlation data, and custom user properties. This is useful for device tracking, message routing, and message lifecycle management scenarios.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {
  default = "mqtt-xokzqqq9"
}

resource "tencentcloud_mqtt_message_enrichment_rule" "example" {
  instance_id = var.instance_id
  rule_name   = "example-rule"
  condition   = base64encode(jsonencode({
    topic= "topic1"
    clientId ="client-1"
    username = "user1" 
  }))
  actions = base64encode(jsonencode({
    responseTopic   = "replies/devices/clientid"
    correlationData = "traceId"
    messageExpiryInterval  = 360
  }))
  status = 1
  remark = "Example message enrichment rule"
  priority = 1
}
```

### With Multiple Conditions

```hcl
resource "tencentcloud_mqtt_message_enrichment_rule" "complex" {
  instance_id = var.instance_id
  rule_name   = "complex-rule"
  condition   = base64encode(jsonencode({
    type = "and"
    rules = [
      {
        field    = "topic"
        operator = "match"
        value    = "device/+/data"
      },
      {
        field    = "qos"
        operator = ">"
        value    = "0"
      }
    ]
  }))
  actions = base64encode(jsonencode([
    {
      type  = "set_property"
      key   = "messageExpiryInterval"
      value = "3600"
    },
    {
      type  = "set_property"
      key   = "responseTopic"
      value = "device/response"
    }
  ]))
  status   = 1
  priority = 10
  remark   = "Complex rule with multiple conditions and actions"
}
```

### Disabled Rule

```hcl
resource "tencentcloud_mqtt_message_enrichment_rule" "disabled" {
  instance_id = var.instance_id
  rule_name   = "disabled-rule"
  condition   = base64encode(jsonencode({
    type = "and"
    rules = [{
      field    = "clientId"
      operator = "match"
      value    = "device-*"
    }]
  }))
  actions = base64encode(jsonencode([{
    type  = "set_property"
    key   = "correlationData"
    value = "tracking-id-123"
  }]))
  status   = 0
  priority = 99
  remark   = "This rule is temporarily disabled"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) MQTT instance ID. You can obtain it from the [DescribeInstanceList](https://cloud.tencent.com/document/api/1778/111029) API or the console.
* `rule_name` - (Required, String) Rule name. Cannot be empty, 3-64 characters, supports Chinese characters, letters, numbers, "-" and "_".
* `condition` - (Required, String) Condition expression. Must be a Base64-encoded JSON string that defines when the rule should be triggered. The JSON structure should contain a `type` field (e.g., "and", "or") and a `rules` array with condition objects.
* `actions` - (Required, String) Actions to perform when the condition is met. Must be a Base64-encoded JSON string array. Each action object should have a `type` field (e.g., "set_property") and relevant parameters like `key` and `value`.
* `status` - (Optional, Int) Rule status. Valid values: `1` for enabled, `0` for disabled. Default: `1`.
* `priority` - (Optional, Int) Priority for rule execution order. Lower values have higher priority. Default: `1`.
* `remark` - (Optional, String) Remark or description for the rule. Maximum 128 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource, in format `{instance_id}#{rule_id}`.
* `rule_id` - The unique rule ID assigned by the cloud service.
* `created_time` - Creation time (Unix timestamp in seconds).
* `update_time` - Last update time (Unix timestamp in seconds).

## Import

MQTT message enrichment rule can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_message_enrichment_rule.example mqtt-xxxxxx#rule-xxxxxx
```

Where:
- `mqtt-xxxxxx` is the instance ID
- `rule-xxxxxx` is the rule ID

## Notes

### Base64 Encoding

Both `condition` and `actions` fields require Base64-encoded JSON strings. In Terraform, you can use the `base64encode()` and `jsonencode()` functions together to properly format these values:

```hcl
condition = base64encode(jsonencode({
  type = "and"
  rules = [...]
}))
```

### Condition Structure

The condition JSON structure typically looks like:

```json
{
  "type": "and",
  "rules": [
    {
      "field": "topic",
      "operator": "==",
      "value": "test/topic"
    }
  ]
}
```

Supported fields: `topic`, `clientId`, `qos`, etc.
Supported operators: `==`, `!=`, `>`, `<`, `>=`, `<=`, `match`, etc.

### Actions Structure

The actions JSON structure is an array of action objects:

```json
[
  {
    "type": "set_property",
    "key": "property_name",
    "value": "property_value"
  }
]
```

Common property keys include:
- `messageExpiryInterval`: Message expiration time (seconds)
- `responseTopic`: Response topic
- `correlationData`: Correlation data
- Custom user properties

### Priority Field

The `priority` field is computed-only and defaults to 1. It cannot be modified by users through Terraform. If you need different priorities for multiple rules, you'll need to manage this through the console or API directly.

### Full Update Semantics

When updating a message enrichment rule, all fields must be provided in the update request. Terraform handles this automatically by reading the current state and submitting all values during updates.

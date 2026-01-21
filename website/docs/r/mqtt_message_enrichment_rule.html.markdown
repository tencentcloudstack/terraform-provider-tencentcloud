---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_message_enrichment_rule"
sidebar_current: "docs-tencentcloud-resource-mqtt_message_enrichment_rule"
description: |-
  Provides a resource to create a MQTT message enrichment rule
---

# tencentcloud_mqtt_message_enrichment_rule

Provides a resource to create a MQTT message enrichment rule

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `actions` - (Required, List) Rule execution actions.
* `condition` - (Required, List) Rule matching condition.
* `instance_id` - (Required, String, ForceNew) MQTT instance ID.
* `priority` - (Required, Int) Rule priority, smaller number means higher priority.
* `rule_name` - (Required, String) Rule name, 3-64 characters, supports Chinese, letters, numbers, `-` and `_`.
* `remark` - (Optional, String) Remark information. not exceeding 128 characters in length.
* `status` - (Optional, Int) Policy status, 0: undefined; 1: active; 2: inactive, default is 2.

The `actions` object supports the following:

* `correlation_data` - (Optional, String) Correlation Data.
* `message_expiry_interval` - (Optional, Int) Message expiration interval.
* `response_topic` - (Optional, String) Response Topic.
* `user_property` - (Optional, List) User Properties.

The `condition` object supports the following:

* `client_id` - (Required, String) Client ID.
* `topic` - (Required, String) Topic.
* `username` - (Required, String) User name.

The `user_property` object of `actions` supports the following:

* `key` - (Required, String) Key.
* `value` - (Required, String) Value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time, millisecond timestamp.
* `rule_id` - Rule ID.
* `update_time` - Update time, millisecond timestamp.


## Import

MQTT message enrichment rule can be imported using the instanceId#ruleId, e.g.

```
terraform import tencentcloud_mqtt_message_enrichment_rule.example mqtt-zxje8zdd#34
```


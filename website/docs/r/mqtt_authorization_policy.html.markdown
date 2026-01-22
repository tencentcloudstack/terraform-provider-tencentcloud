---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_authorization_policy"
sidebar_current: "docs-tencentcloud-resource-mqtt_authorization_policy"
description: |-
  Provides a resource to create a MQTT authorization policy
---

# tencentcloud_mqtt_authorization_policy

Provides a resource to create a MQTT authorization policy

## Example Usage

```hcl
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
```

### Or

```hcl
resource "tencentcloud_mqtt_authorization_policy" "example" {
  instance_id    = "mqtt-g4qgr3gx"
  policy_name    = "tf-example"
  policy_version = 1
  priority       = 10
  effect         = "deny"
  actions        = "pub,sub"
  retain         = 3
  qos            = "1,2"
  resources      = "topic-demo"
  username       = "root*"
  client_id      = "*$${Username}*"
  ip             = "192.168.1.0/24"
  remark         = "policy remark."
}
```

## Argument Reference

The following arguments are supported:

* `actions` - (Required, String) Operation - connect: connect; pub: publish; sub: subscribe.
* `effect` - (Required, String) Decision: allow/deny.
* `instance_id` - (Required, String, ForceNew) MQTT instance ID.
* `policy_name` - (Required, String) Policy name, cannot be empty, 3-64 characters, supports Chinese characters, letters, numbers, "-" and "_".
* `policy_version` - (Required, Int) Policy version, default is 1, currently only 1 is supported.
* `priority` - (Required, Int) The strategy priority, the smaller the higher the priority, cannot be repeated.
* `qos` - (Required, String) Condition: Quality of Service 0: At most once 1: At least once 2: Exactly once.
* `retain` - (Required, Int) Condition - Reserved message 1, match reserved message; 2, match unreserved message, 3. match reserved and unreserved message.
* `client_id` - (Optional, String) Condition - Client ID, supports regular expressions.
* `ip` - (Optional, String) Condition - Client IP address, supports IP or CIDR.
* `remark` - (Optional, String) Remarks, up to 128 characters.
* `resources` - (Optional, String) Resources, requiring matching subscriptions.
* `username` - (Optional, String) Condition - Username.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `policy_id` - Authorization policy rule id.


## Import

MQTT authorization policy can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_authorization_policy.example mqtt-g4qgr3gx#140
```


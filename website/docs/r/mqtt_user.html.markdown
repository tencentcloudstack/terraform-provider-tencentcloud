---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_user"
sidebar_current: "docs-tencentcloud-resource-mqtt_user"
description: |-
  Provides a resource to create a MQTT user
---

# tencentcloud_mqtt_user

Provides a resource to create a MQTT user

## Example Usage

```hcl
resource "tencentcloud_mqtt_user" "example" {
  instance_id = "mqtt-zxjwkr98"
  username    = "tf-example"
  password    = "Password@123"
  remark      = "Remark."
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `username` - (Required, String, ForceNew) Username, cannot be empty, only supports uppercase and lowercase letter separators ("_", "-"), cannot exceed 32 characters.
* `password` - (Optional, String) Password, when this field is empty, the backend will generate it by default.
* `remark` - (Optional, String) Note that the length should not exceed 128 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Creation time, millisecond timestamp.
* `modified_time` - Modify time, millisecond timestamp.


## Import

MQTT user can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_user.example mqtt-zxjwkr98#tf-example
```


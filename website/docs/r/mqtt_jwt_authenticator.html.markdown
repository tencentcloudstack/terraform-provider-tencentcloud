---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_jwt_authenticator"
sidebar_current: "docs-tencentcloud-resource-mqtt_jwt_authenticator"
description: |-
  Provides a resource to create a MQTT jwt authenticator
---

# tencentcloud_mqtt_jwt_authenticator

Provides a resource to create a MQTT jwt authenticator

## Example Usage

### If algorithm is hmac-based

```hcl
resource "tencentcloud_mqtt_jwt_authenticator" "example" {
  instance_id = "mqtt-zxjwkr98"
  algorithm   = "hmac-based"
  from        = "password"
  secret      = "your secret content"
  remark      = "Remark."
}
```

### If algorithm is public-key

```hcl
resource "tencentcloud_mqtt_jwt_authenticator" "example" {
  instance_id = "mqtt-zxjwkr98"
  algorithm   = "public-key"
  from        = "username"
  public_key  = "your public key"
  remark      = "Remark."
}
```

## Argument Reference

The following arguments are supported:

* `algorithm` - (Required, String) Algorithm. hmac-based, public-key.
* `from` - (Required, String) Pass the key of JWT when connecting the device; Username - passed using the username field; Password - Pass using password field.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `public_key` - (Optional, String) Public key.
* `remark` - (Optional, String) Remark.
* `secret` - (Optional, String) Secret.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

MQTT jwt authenticator can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_jwt_authenticator.example mqtt-zxjwkr98
```


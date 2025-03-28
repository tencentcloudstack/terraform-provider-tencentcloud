---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_http_authenticator"
sidebar_current: "docs-tencentcloud-resource-mqtt_http_authenticator"
description: |-
  Provides a resource to create a MQTT http authenticator
---

# tencentcloud_mqtt_http_authenticator

Provides a resource to create a MQTT http authenticator

## Example Usage

```hcl
resource "tencentcloud_mqtt_http_authenticator" "example" {
  instance_id     = "mqtt-zxjwkr98"
  endpoint        = "https://example.com"
  concurrency     = 8
  method          = "POST"
  status          = "open"
  remark          = "Remark."
  connect_timeout = 10
  read_timeout    = 10
  header {
    key   = "Content-type"
    value = "application/json"
  }

  body {
    key   = "bodyKey"
    value = "bodyValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `endpoint` - (Required, String) JWKS endpoint.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `body` - (Optional, List) Forwarding request body.
* `concurrency` - (Optional, Int) Maximum concurrent connections, default 8, range: 1-20.
* `connect_timeout` - (Optional, Int) Connection timeout, unit: seconds, range: 1-30.
* `header` - (Optional, List) Forwarding request header.
* `method` - (Optional, String) Network request method GET or POST, default POST.
* `read_timeout` - (Optional, Int) Request timeout, unit: seconds, range: 1-30.
* `remark` - (Optional, String) Remark.
* `status` - (Optional, String) Is the authenticator enabled: open enable; Close close.

The `body` object supports the following:

* `key` - (Required, String) Body key.
* `value` - (Required, String) Body key.

The `header` object supports the following:

* `key` - (Required, String) Header key.
* `value` - (Required, String) Header value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

MQTT http authenticator can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_http_authenticator.example mqtt-zxjwkr98
```


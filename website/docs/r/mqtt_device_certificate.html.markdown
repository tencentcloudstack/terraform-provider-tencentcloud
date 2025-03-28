---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_device_certificate"
sidebar_current: "docs-tencentcloud-resource-mqtt_device_certificate"
description: |-
  Provides a resource to create a MQTT device certificate
---

# tencentcloud_mqtt_device_certificate

Provides a resource to create a MQTT device certificate

## Example Usage

```hcl
resource "tencentcloud_mqtt_device_certificate" "example" {
  instance_id        = "mqtt-zxjwkr98"
  device_certificate = ""
  ca_sn              = ""
  client_id          = ""
  format             = ""
  status             = "ACTIVE"
}
```

## Argument Reference

The following arguments are supported:

* `ca_sn` - (Required, String) Associated CA certificate SN.
* `device_certificate` - (Required, String) Device certificate.
* `instance_id` - (Required, String) Instance ID.
* `client_id` - (Optional, String) Client ID.
* `format` - (Optional, String) Certificate format, Default is PEM.
* `status` - (Optional, String) Certificate status, Default is ACTIVE.\n  ACTIVE activation;\n  INACTIVE not active.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `certificate_source` - Certificate source.
* `created_time` - Certificate create time.
* `device_certificate_cn` - Certificate common name.
* `device_certificate_sn` - Equipment certificate serial number.
* `not_after_time` - Certificate expiring date.
* `not_before_time` - Certificate effective start date.
* `update_time` - Certificate update time.


## Import

MQTT device certificate can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_device_certificate.example mqtt_device_certificate_id
```


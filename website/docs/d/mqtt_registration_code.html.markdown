---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_registration_code"
sidebar_current: "docs-tencentcloud-datasource-mqtt_registration_code"
description: |-
  Use this data source to query detailed information of MQTT registration code
---

# tencentcloud_mqtt_registration_code

Use this data source to query detailed information of MQTT registration code

## Example Usage

```hcl
data "tencentcloud_mqtt_registration_code" "example" {
  instance_id = "mqtt-zxjwkr98"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `registration_code` - Registration code.



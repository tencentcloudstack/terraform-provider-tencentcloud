---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_topics"
sidebar_current: "docs-tencentcloud-datasource-mqtt_topics"
description: |-
  Use this data source to query detailed information of MQTT topics
---

# tencentcloud_mqtt_topics

Use this data source to query detailed information of MQTT topics

## Example Usage

```hcl
data "tencentcloud_mqtt_topics" "example" {
  instance_id = "mqtt-g4qgr3gx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Topic list.



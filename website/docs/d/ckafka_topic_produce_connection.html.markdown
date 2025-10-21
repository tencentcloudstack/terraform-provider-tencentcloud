---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_topic_produce_connection"
sidebar_current: "docs-tencentcloud-datasource-ckafka_topic_produce_connection"
description: |-
  Use this data source to query detailed information of ckafka topic_produce_connection
---

# tencentcloud_ckafka_topic_produce_connection

Use this data source to query detailed information of ckafka topic_produce_connection

## Example Usage

```hcl
data "tencentcloud_ckafka_topic_produce_connection" "topic_produce_connection" {
  instance_id = "ckafka-xxxxxx"
  topic_name  = "topic-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) InstanceId.
* `topic_name` - (Required, String) TopicName.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - link information return result set.
  * `ip_addr` - ip address.
  * `is_un_support_version` - Is the supported version.
  * `time` - connect time.



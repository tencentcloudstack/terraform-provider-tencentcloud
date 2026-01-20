---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_version"
sidebar_current: "docs-tencentcloud-datasource-ckafka_version"
description: |-
  Use this data source to query detailed information of CKafka version
---

# tencentcloud_ckafka_version

Use this data source to query detailed information of CKafka version

## Example Usage

```hcl
data "tencentcloud_ckafka_version" "example" {
  instance_id = "ckafka-8j4raxv8"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) CKafka instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cur_broker_version` - Current broker version.
* `kafka_version` - Current Kafka version.
* `latest_broker_versions` - List of latest broker versions supported by the platform.
  * `broker_version` - Broker version.
  * `kafka_version` - Kafka version.



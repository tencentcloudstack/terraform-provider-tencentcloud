---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_publisher_summary"
sidebar_current: "docs-tencentcloud-datasource-tdmq_publisher_summary"
description: |-
  Use this data source to query detailed information of tdmq publisher_summary
---

# tencentcloud_tdmq_publisher_summary

Use this data source to query detailed information of tdmq publisher_summary

## Example Usage

```hcl
data "tencentcloud_tdmq_publisher_summary" "publisher_summary" {
  cluster_id = "pulsar-9n95ax58b9vn"
  namespace  = "keep-ns"
  topic      = "keep-topic"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `namespace` - (Required, String) namespace name.
* `topic` - (Required, String) subject name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `msg_rate_in` - Production rate (units per second)Note: This field may return null, indicating that no valid value can be obtained.
* `msg_throughput_in` - Production rate (bytes per second)Note: This field may return null, indicating that no valid value can be obtained.
* `publisher_count` - number of producersNote: This field may return null, indicating that no valid value can be obtained.
* `storage_size` - Message store size in bytesNote: This field may return null, indicating that no valid value can be obtained.



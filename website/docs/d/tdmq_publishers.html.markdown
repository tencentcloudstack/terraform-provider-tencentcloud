---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_publishers"
sidebar_current: "docs-tencentcloud-datasource-tdmq_publishers"
description: |-
  Use this data source to query detailed information of tdmq publishers
---

# tencentcloud_tdmq_publishers

Use this data source to query detailed information of tdmq publishers

## Example Usage

```hcl
data "tencentcloud_tdmq_publishers" "publishers" {
  cluster_id = "pulsar-9n95ax58b9vn"
  namespace  = "keep-ns"
  topic      = "keep-topic"
  filters {
    name   = "ProducerName"
    values = ["test"]
  }
  sort {
    name  = "ProducerName"
    order = "DESC"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `namespace` - (Required, String) namespace name.
* `topic` - (Required, String) topic name.
* `filters` - (Optional, List) Parameter filter, support ProducerName, Address field.
* `result_output_file` - (Optional, String) Used to save results.
* `sort` - (Optional, List) sorter.

The `filters` object supports the following:

* `name` - (Optional, String) The name of the filter parameter.
* `values` - (Optional, Set) value.

The `sort` object supports the following:

* `name` - (Required, String) sorter.
* `order` - (Required, String) Ascending ASC, descending DESC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `publishers` - Producer Information ListNote: This field may return null, indicating that no valid value can be obtained.
  * `address` - producer addressNote: This field may return null, indicating that no valid value can be obtained.
  * `average_msg_size` - Average message size (bytes)Note: This field may return null, indicating that no valid value can be obtained.
  * `client_version` - client versionNote: This field may return null, indicating that no valid value can be obtained.
  * `connected_since` - connection timeNote: This field may return null, indicating that no valid value can be obtained.
  * `msg_rate_in` - Message production rate (articles/second)Note: This field may return null, indicating that no valid value can be obtained.
  * `msg_throughput_in` - Message production throughput rate (bytes/second)Note: This field may return null, indicating that no valid value can be obtained.
  * `partition` - The topic partition number of the producer connectionNote: This field may return null, indicating that no valid value can be obtained.
  * `producer_id` - producer idNote: This field may return null, indicating that no valid value can be obtained.
  * `producer_name` - producer nameNote: This field may return null, indicating that no valid value can be obtained.



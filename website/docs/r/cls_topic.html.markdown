---
subcategory: "CLS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_topic"
sidebar_current: "docs-tencentcloud-resource-cls_topic"
description: |-
  Provides a resource to create a cls topic.
---

# tencentcloud_cls_topic

Provides a resource to create a cls topic.

## Example Usage

```hcl
resource "tencentcloud_cls_topic" "topic" {
  topic_name           = "topic"
  logset_id            = "5cd3a17e-fb0b-418c-afd7-77b365397426"
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags = {
    "test" = "test",
  }
}
```

## Argument Reference

The following arguments are supported:

* `logset_id` - (Required) Logset ID.
* `topic_name` - (Required) Log topic name.
* `auto_split` - (Optional) Whether to enable automatic split. Default value: true.
* `max_split_partitions` - (Optional) Maximum number of partitions to split into for this topic if automatic split is enabled. Default value: 50.
* `partition_count` - (Optional, ForceNew) Number of log topic partitions. Default value: 1. Maximum value: 10.
* `period` - (Optional) Lifecycle in days. Value range: 1~366. Default value: 30.
* `storage_type` - (Optional, ForceNew) Log topic storage class. Valid values: hot: real-time storage; cold: offline storage. Default value: hot. If cold is passed in, please contact the customer service to add the log topic to the allowlist first..
* `tags` - (Optional) Tag description list. Up to 10 tag key-value pairs are supported and must be unique.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls topic can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_topic.topic 2f5764c1-c833-44c5-84c7-950979b2a278
```


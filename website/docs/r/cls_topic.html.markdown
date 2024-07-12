---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_topic"
sidebar_current: "docs-tencentcloud-resource-cls_topic"
description: |-
  Provides a resource to create a cls topic.
---

# tencentcloud_cls_topic

Provides a resource to create a cls topic.

## Example Usage

### Create a standard cls topic

```hcl
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf_example"
  tags = {
    tagKey = "tagValue"
  }
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf_example"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 30
  storage_type         = "hot"
  describes            = "Test Demo."
  hot_period           = 10
  tags = {
    tagKey = "tagValue"
  }
}
```

### Create a cls topic with web tracking

```hcl
resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf_example"
  tags = {
    tagKey = "tagValue"
  }
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf_example"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 30
  storage_type         = "hot"
  describes            = "Test Demo."
  hot_period           = 10
  is_web_tracking      = true

  extends {
    anonymous_access {
      operations = ["trackLog", "realtimeProducer"]
      conditions {
        attributes      = "VpcID"
        rule            = 1
        condition_value = "vpc-ahr3xajx"
      }
    }
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `logset_id` - (Required, String) Logset ID.
* `topic_name` - (Required, String) Log topic name.
* `auto_split` - (Optional, Bool) Whether to enable automatic split. Default value: true.
* `describes` - (Optional, String) Log Topic Description.
* `extends` - (Optional, List) Log Subject Extension Information.
* `hot_period` - (Optional, Int) 0: Turn off log sinking. Non 0: The number of days of standard storage after enabling log settling. HotPeriod needs to be greater than or equal to 7 and less than Period. Only effective when StorageType is hot.
* `is_web_tracking` - (Optional, Bool) No authentication switch. False: closed; True: Enable. The default is false. After activation, anonymous access to the log topic will be supported for specified operations.
* `max_split_partitions` - (Optional, Int) Maximum number of partitions to split into for this topic if automatic split is enabled. Default value: 50.
* `partition_count` - (Optional, Int) Number of log topic partitions. Default value: 1. Maximum value: 10.
* `period` - (Optional, Int) Lifecycle in days. Value range: 1~366. Default value: 30.
* `storage_type` - (Optional, String) Log topic storage class. Valid values: hot: real-time storage; cold: offline storage. Default value: hot. If cold is passed in, please contact the customer service to add the log topic to the allowlist first.
* `tags` - (Optional, Map) Tag description list. Up to 10 tag key-value pairs are supported and must be unique.

The `anonymous_access` object of `extends` supports the following:

* `conditions` - (Optional, List) Operation list, supporting trackLog (JS/HTTP upload log) and realtimeProducer (kafka protocol upload log).
* `operations` - (Optional, List) Operation list, supporting trackLog (JS/HTTP upload log) and realtimeProducer (kafka protocol upload log).

The `conditions` object of `anonymous_access` supports the following:

* `attributes` - (Optional, String) Condition attribute, currently only VpcID is supported.
* `condition_value` - (Optional, String) The value of the corresponding conditional attribute.
* `rule` - (Optional, Int) Conditional rule, 1: equal, 2: not equal.

The `extends` object supports the following:

* `anonymous_access` - (Optional, List) Log topic authentication free configuration information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls topic can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_topic.example 2f5764c1-c833-44c5-84c7-950979b2a278
```


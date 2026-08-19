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

~> **NOTE:** Field `encryption` can only be enabled, not disabled.

~> **NOTE:** Field `custom_kms_info` is for user-defined KMS key. If not set, the CLS default key (alias KMS-CLS) is used.

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

### Create a cls metric topic(biz_type=1)

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
  biz_type             = 1
  tags = {
    tagKey = "tagValue"
  }
}
```

### Create a cls topic with custom KMS key (encryption=1)

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
  encryption           = 1

  custom_kms_info {
    kms_region = "ap-guangzhou"
    kms_key_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `logset_id` - (Required, String) Logset ID. Get the logset ID via `DescribeLogsets` API.
* `topic_name` - (Required, String) Log topic name. Constraints: cannot be an empty string, cannot contain the `|` character, and cannot use the following reserved names: `cls_service_log`, `loglistener_status`, `loglistener_alarm`, `loglistener_business`, `cls_service_metric`.
* `auto_split` - (Optional, Bool) Whether to enable automatic split. Default value: `true`.
* `biz_type` - (Optional, Int) Topic type. `0`: log topic (default), `1`: metric topic.
* `custom_kms_info` - (Optional, List) User-defined KMS key information. If empty, the default key (alias `KMS-CLS`) is used.
* `describes` - (Optional, String) Log topic description.
* `encryption` - (Optional, Int) Encryption-related parameters. Supported for encryption-enabled regions and allowlisted users; cannot be passed in other scenarios. `0` or not passed: no encryption; `1`: kms-cls cloud product key encryption. Once enabled, it cannot be disabled. Supported regions: ap-beijing, ap-guangzhou, ap-shanghai, ap-singapore, ap-bangkok, ap-jakarta, eu-frankfurt, ap-seoul, ap-tokyo.
* `extends` - (Optional, List) Topic extension information.
* `hot_period` - (Optional, Int) `0`: turn off log settling. Non-`0`: the number of days of standard storage after enabling log settling. HotPeriod must be greater than or equal to 7 and less than Period. Only effective when `storage_type` is `hot`. Not supported for metric topics.
* `is_web_tracking` - (Optional, Bool) Free authentication switch. `false`: closed (default); `true`: enabled. When enabled, anonymous access to the log topic will be supported for specified operations. Not supported for metric topics.
* `max_split_partitions` - (Optional, Int) Maximum number of partitions allowed for the topic if automatic split is enabled. Default value: `50`.
* `partition_count` - (Optional, Int) Number of log topic partitions. Default: 1, maximum: 10.
* `period` - (Optional, Int) Retention period, unit: days. Log topic (standard storage): 1 to 3600 days, value `3640` means permanent retention. Log topic (infrequent storage): 7 to 3600 days, value `3640` means permanent retention. Metric topic: 1 to 3600 days, value `3640` means permanent retention.
* `storage_type` - (Optional, String) Log topic storage type. Valid values: `hot`: standard storage; `cold`: infrequent storage. Default value: `hot`. Not supported for metric topics.
* `tags` - (Optional, Map) Tag description list. Up to 10 tag key-value pairs are supported, and the same resource can only be bound to the same tag key.

The `anonymous_access` object of `extends` supports the following:

* `conditions` - (Optional, List) Operation list, supporting trackLog (JS/HTTP upload log) and realtimeProducer (kafka protocol upload log).
* `operations` - (Optional, List) Operation list, supporting trackLog (JS/HTTP upload log) and realtimeProducer (kafka protocol upload log).

The `conditions` object of `anonymous_access` supports the following:

* `attributes` - (Optional, String) Condition attribute, currently only VpcID is supported.
* `condition_value` - (Optional, String) The value of the corresponding conditional attribute.
* `rule` - (Optional, Int) Conditional rule, 1: equal, 2: not equal.

The `custom_kms_info` object supports the following:

* `kms_key_id` - (Required, String) KMS key ID.
* `kms_region` - (Required, String) KMS region. Refer to Tencent Cloud KMS documentation for supported regions. Format: `ap-guangzhou`.

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

cls metric topic (biz_type=1) can be imported using the id with "#1" suffix, e.g.

```
$ terraform import tencentcloud_cls_topic.example 2f5764c1-c833-44c5-84c7-950979b2a278#1
```


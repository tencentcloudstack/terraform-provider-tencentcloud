---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_consumer_group_modify_offset"
sidebar_current: "docs-tencentcloud-resource-ckafka_consumer_group_modify_offset"
description: |-
  Provides a resource to create a ckafka consumer_group_modify_offset
---

# tencentcloud_ckafka_consumer_group_modify_offset

Provides a resource to create a ckafka consumer_group_modify_offset

## Example Usage

```hcl
resource "tencentcloud_ckafka_consumer_group_modify_offset" "consumer_group_modify_offset" {
  instance_id = "ckafka-xxxxxx"
  group       = "xxxxxx"
  offset      = 0
  strategy    = 2
  topics      = ["xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required, String, ForceNew) kafka group.
* `instance_id` - (Required, String, ForceNew) Kafka instance id.
* `strategy` - (Required, Int, ForceNew) Reset the policy of offset.
`0`: Move the offset forward or backward shift bar;
`1`: Alignment reference (by-duration,to-datetime,to-earliest,to-latest), which means moving the offset to the location of the specified timestamp;
`2`: Alignment reference (to-offset), which means to move the offset to the specified offset location.
* `offset` - (Optional, Int, ForceNew) The offset location that needs to be reset. When strategy is 2, this field must be included.
* `partitions` - (Optional, Set: [`Int`], ForceNew) The list of partition that needs to be reset if no Topics parameter is specified. Resets the partition in the corresponding Partition list of all topics. When Topics is specified, the partition of the corresponding topic list of the specified Partitions list is reset.
* `shift_timestamp` - (Optional, Int, ForceNew) Unit ms. When strategy is 1, you must include this field, where-2 means to reset the offset to the beginning,-1 means to reset to the latest position (equivalent to emptying), and other values represent the specified time. You will get the offset of the specified time in the topic and then reset it. If there is no message at the specified time, get the last offset.
* `shift` - (Optional, Int, ForceNew) This field must be included when strategy is 0. If it is greater than zero, the offset will be moved backward by shift bars, and if it is less than zero, the offset will be traced back to the number of shift entries. After the correct reset, the new offset should be (old_offset + shift). It should be noted that if the new offset is less than partition's earliest, it will be set to earliest, and if the latest greater than partition will be set to latest.
* `topics` - (Optional, Set: [`String`], ForceNew) Indicates the topics that needs to be reset. Leave it empty means all.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.




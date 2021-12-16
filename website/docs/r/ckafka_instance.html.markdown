---
subcategory: "Ckafka"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_instance"
sidebar_current: "docs-tencentcloud-resource-ckafka_instance"
description: |-
  Use this resource to create ckafka instance.
---

# tencentcloud_ckafka_instance

Use this resource to create ckafka instance.

## Example Usage

```hcl
resource "tencentcloud_ckafka_instance" "foo" {
  instance_name      = "demo-hello"
  zone_id            = 100006
  period             = 1
  vpc_id             = "vpc-boi1ah65"
  subnet_id          = "subnet-7ros461e"
  msg_retention_time = 1440

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    enable = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) Instance name.
* `period` - (Required, ForceNew) Prepaid purchase time, such as 1, is one month.
* `subnet_id` - (Required, ForceNew) Subnet id.
* `vpc_id` - (Required, ForceNew) Vpc id.
* `zone_id` - (Required, ForceNew) Available zone id.
* `band_width` - (Optional, ForceNew) Whether to open the ip whitelist, `true`: open, `false`: close.
* `config` - (Optional) Instance configuration.
* `disk_size` - (Optional, ForceNew) Disk Size.
* `disk_type` - (Optional, ForceNew) Type of disk.
* `dynamic_retention_config` - (Optional) Dynamic message retention policy configuration.
* `kafka_version` - (Optional, ForceNew) Kafka version (0.10.2/1.1.1/2.4.1).
* `msg_retention_time` - (Optional) The maximum retention time of instance logs, in minutes. the default is 10080 (7 days), the maximum is 30 days, and the default 0 is not filled, which means that the log retention time recovery policy is not enabled.
* `partition` - (Optional, ForceNew) Partition size, the professional version does not need set.
* `public_network` - (Optional) Timestamp.
* `rebalance_time` - (Optional) Modification of the rebalancing time after upgrade.
* `renew_flag` - (Optional, ForceNew) Prepaid automatic renewal mark, 0 means the default state, the initial state, 1 means automatic renewal, 2 means clear no automatic renewal (user setting).
* `tags` - (Optional, ForceNew) Partition size, the professional version does not need tag.

The `config` object supports the following:

* `auto_create_topic_enable` - (Required) Automatic creation. true: enabled, false: not enabled.
* `default_num_partitions` - (Required) If auto.create.topic.enable is set to true and this value is not set, 3 will be used by default.
* `default_replication_factor` - (Required) If auto.create.topic.enable is set to true but this value is not set, 2 will be used by default.

The `dynamic_retention_config` object supports the following:

* `bottom_retention` - (Optional) Minimum retention time, in minutes.
* `disk_quota_percentage` - (Optional) Disk quota threshold (in percentage) for triggering the message retention time change event.
* `enable` - (Optional) Whether the dynamic message retention time configuration is enabled. 0: disabled; 1: enabled.
* `step_forward_percentage` - (Optional) Percentage by which the message retention time is shortened each time.

The `tags` object supports the following:

* `key` - (Required) Tag key.
* `value` - (Required) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ckafka instance can be imported using the instance_id, e.g.

```
$ terraform import tencentcloud_ckafka_instance.foo ckafka-f9ife4zz
```


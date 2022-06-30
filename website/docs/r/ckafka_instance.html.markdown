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

~> **NOTE:** It only support create profession ckafka instance.

## Example Usage

```hcl
resource "tencentcloud_ckafka_instance" "foo" {
  band_width         = 40
  disk_size          = 500
  disk_type          = "CLOUD_BASIC"
  period             = 1
  instance_name      = "ckafka-instance-tf-test"
  kafka_version      = "1.1.1"
  msg_retention_time = 1300
  multi_zone_flag    = true
  partition          = 800
  public_network     = 3
  renew_flag         = 0
  subnet_id          = "subnet-4vwihrzk"
  vpc_id             = "vpc-82p1t1nv"
  zone_id            = 100006
  zone_ids = [
    100006,
    100007,
  ]

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    bottom_retention        = 0
    disk_quota_percentage   = 0
    enable                  = 1
    step_forward_percentage = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) Instance name.
* `subnet_id` - (Required) Subnet id.
* `vpc_id` - (Required) Vpc id.
* `zone_id` - (Required) Available zone id.
* `band_width` - (Optional) Instance bandwidth in MBps.
* `config` - (Optional) Instance configuration.
* `disk_size` - (Optional) Disk Size. Its interval varies with bandwidth, and the input must be within the interval, which can be viewed through the control. If it is not within the interval, the plan will cause a change when first created.
* `disk_type` - (Optional) Type of disk.
* `dynamic_retention_config` - (Optional) Dynamic message retention policy configuration.
* `kafka_version` - (Optional) Kafka version (0.10.2/1.1.1/2.4.1).
* `msg_retention_time` - (Optional) The maximum retention time of instance logs, in minutes. the default is 10080 (7 days), the maximum is 30 days, and the default 0 is not filled, which means that the log retention time recovery policy is not enabled.
* `multi_zone_flag` - (Optional) Indicates whether the instance is multi zones. NOTE: if set to `true`, `zone_ids` must set together.
* `partition` - (Optional) Partition Size. Its interval varies with bandwidth, and the input must be within the interval, which can be viewed through the control. If it is not within the interval, the plan will cause a change when first created.
* `period` - (Optional) Prepaid purchase time, such as 1, is one month.
* `public_network` - (Optional) Timestamp.
* `rebalance_time` - (Optional) Modification of the rebalancing time after upgrade.
* `renew_flag` - (Optional) Prepaid automatic renewal mark, 0 means the default state, the initial state, 1 means automatic renewal, 2 means clear no automatic renewal (user setting).
* `tags` - (Optional) Partition size, the professional version does not need tag.
* `zone_ids` - (Optional) List of available zone id. NOTE: this argument must set together with `multi_zone_flag`.

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
* `vip` - Vip of instance.
* `vport` - Type of instance.


## Import

ckafka instance can be imported using the instance_id, e.g.

```
$ terraform import tencentcloud_ckafka_instance.foo ckafka-f9ife4zz
```


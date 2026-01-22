---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_instance"
sidebar_current: "docs-tencentcloud-resource-ckafka_instance"
description: |-
  Use this resource to create CKafka instance.
---

# tencentcloud_ckafka_instance

Use this resource to create CKafka instance.

## Example Usage

### Create basic instance(prepaid)

```hcl
data "tencentcloud_availability_zones_by_product" "gz" {
  name    = "ap-guangzhou-6"
  product = "ckafka"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create ckafka
resource "tencentcloud_ckafka_instance" "example" {
  instance_name       = "tf-example"
  zone_id             = data.tencentcloud_availability_zones_by_product.gz.zones.0.id
  period              = 1
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  msg_retention_time  = 1300
  renew_flag          = 0
  kafka_version       = "2.8.1"
  disk_size           = 200
  disk_type           = "CLOUD_BASIC"
  band_width          = 40
  partition           = 400
  specifications_type = "profession"
  instance_type       = 1

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

### Create multi zone instance(postpaid)

```hcl
data "tencentcloud_availability_zones_by_product" "gz6" {
  name    = "ap-guangzhou-6"
  product = "ckafka"
}

data "tencentcloud_availability_zones_by_product" "gz7" {
  name    = "ap-guangzhou-7"
  product = "ckafka"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create ckafka
resource "tencentcloud_ckafka_instance" "example" {
  instance_name   = "tf-example"
  zone_id         = data.tencentcloud_availability_zones_by_product.gz6.zones.0.id
  multi_zone_flag = true
  zone_ids = [
    data.tencentcloud_availability_zones_by_product.gz6.zones.0.id,
    data.tencentcloud_availability_zones_by_product.gz7.zones.0.id,
  ]
  renew_flag          = 0
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  msg_retention_time  = 4320
  kafka_version       = "2.8.1"
  disk_size           = 200
  disk_type           = "CLOUD_BASIC"
  band_width          = 20
  partition           = 400
  specifications_type = "profession"
  charge_type         = "POSTPAID_BY_HOUR"
  instance_type       = 1

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

* `instance_name` - (Required, String) Instance name.
* `zone_id` - (Required, Int) Available zone id.
* `band_width` - (Optional, Int) Instance bandwidth in MBps.
* `charge_type` - (Optional, String, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `PREPAID`.
* `config` - (Optional, List) Instance configuration.
* `custom_ssl_cert_id` - (Optional, String) Custom certificate ID, only effective when `specifications_type` is set to `profession`, supports custom certificate capabilities.
* `disk_size` - (Optional, Int) Disk Size. Its interval varies with bandwidth, and the input must be within the interval, which can be viewed through the control. If it is not within the interval, the plan will cause a change when first created.
* `disk_type` - (Optional, String) Type of disk.
* `dynamic_retention_config` - (Optional, List) Dynamic message retention policy configuration.
* `elastic_bandwidth_switch` - (Optional, Int) Elastic bandwidth switch 0 not turned on 1 turned on (0 default). This takes effect only when the instance is created.
* `instance_type` - (Optional, Int) Description of instance type. `profession`: 1, `standard`:  1(general), 2(standard), 3(advanced), 4(capacity), 5(specialized-1), 6(specialized-2), 7(specialized-3), 8(specialized-4), 9(exclusive).
* `kafka_version` - (Optional, String) Kafka version (0.10.2/1.1.1/2.4.1).
* `max_message_byte` - (Optional, Int) The size of a single message in bytes at the instance level. Value range: `1024 - 12*1024*1024 bytes (i.e., 1KB-12MB).
* `msg_retention_time` - (Optional, Int) The maximum retention time of instance logs, in minutes. the default is 10080 (7 days), the maximum is 30 days, and the default 0 is not filled, which means that the log retention time recovery policy is not enabled.
* `multi_zone_flag` - (Optional, Bool) Indicates whether the instance is multi zones. NOTE: if set to `true`, `zone_ids` must set together.
* `partition` - (Optional, Int) Partition Size. Its interval varies with bandwidth, and the input must be within the interval, which can be viewed through the control. If it is not within the interval, the plan will cause a change when first created.
* `period` - (Optional, Int) Prepaid purchase time, such as 1, is one month.
* `public_network` - (Optional, Int) Bandwidth of the public network.
* `rebalance_time` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.82.37. Modification of the rebalancing time after upgrade.
* `renew_flag` - (Optional, Int) Prepaid automatic renewal mark, 0 means the default state, the initial state, 1 means automatic renewal, 2 means clear no automatic renewal (user setting).
* `specifications_type` - (Optional, String) Specifications type of instance. Allowed values are `profession`, `premium`. Default is `profession`.
* `subnet_id` - (Optional, String) Subnet id, it will be basic network if not set.
* `tag_set` - (Optional, Map) Tag set of instance.
* `tags` - (Optional, List, **Deprecated**) It has been deprecated from version 1.78.5, because it do not support change. Use `tag_set` instead. Tags of instance. Partition size, the professional version does not need tag.
* `upgrade_strategy` - (Optional, Int) POSTPAID_BY_HOUR scale-down mode
- 1: stable transformation;
- 2: High-speed transformer.
* `vpc_id` - (Optional, String) Vpc id, it will be basic network if not set.
* `zone_ids` - (Optional, Set: [`Int`]) List of available zone id. NOTE: this argument must set together with `multi_zone_flag`.

The `config` object supports the following:

* `auto_create_topic_enable` - (Required, Bool) Automatic creation. true: enabled, false: not enabled.
* `default_num_partitions` - (Required, Int) If auto.create.topic.enable is set to true and this value is not set, 3 will be used by default.
* `default_replication_factor` - (Required, Int) If auto.create.topic.enable is set to true but this value is not set, 2 will be used by default.

The `dynamic_retention_config` object supports the following:

* `bottom_retention` - (Optional, Int) Minimum retention time, in minutes.
* `disk_quota_percentage` - (Optional, Int) Disk quota threshold (in percentage) for triggering the message retention time change event.
* `enable` - (Optional, Int) Whether the dynamic message retention time configuration is enabled. 0: disabled; 1: enabled.
* `step_forward_percentage` - (Optional, Int) Percentage by which the message retention time is shortened each time.

The `tags` object supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `vip` - Vip of instance.
* `vport` - Type of instance.


## Import

CKafka instance can be imported using the instanceId, e.g.

```
$ terraform import tencentcloud_ckafka_instance.example ckafka-f9ife4zz
```


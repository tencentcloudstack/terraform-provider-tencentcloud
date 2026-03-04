---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_instances_v2"
sidebar_current: "docs-tencentcloud-datasource-ckafka_instances_v2"
description: |-
  Use this data source to query detailed instance information of Ckafka instances
---

# tencentcloud_ckafka_instances_v2

Use this data source to query detailed instance information of Ckafka instances

## Example Usage

### Query all Ckafka instances

```hcl
data "tencentcloud_ckafka_instances_v2" "example" {}
```

### Query Ckafka instances by filters

```hcl
data "tencentcloud_ckafka_instances_v2" "example" {
  filters {
    name   = "InstanceId"
    values = ["ckafka-7k5nbnem"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions. Supported fields: Ip, VpcId, SubNetId, InstanceType, InstanceId. Note: filter.Values can contain up to 10 values.
* `instance_id_list` - (Optional, Set: [`String`]) Filter by instance ID list.
* `instance_id` - (Optional, String) Filter by instance ID.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) Filter by instance name, instance ID, availability zone, VPC ID or subnet ID. Fuzzy search is supported.
* `status` - (Optional, Set: [`Int`]) Filter by instance status. -1: creation failed, 0: creating, 1: running, 2: deleting, 3: deleted, 4: deletion failed, 5: isolated, 7: upgrading.
* `tag_key` - (Optional, String) Filter by tag key.
* `tag_list` - (Optional, List) Filter by tag list (intersection).

The `filters` object supports the following:

* `name` - (Required, String) Filter field name. Supported: Ip, VpcId, SubNetId, InstanceType, InstanceId.
* `values` - (Required, Set) Filter field values (up to 10 values).

The `tag_list` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - Instance list.
  * `bandwidth` - Instance bandwidth in Mbps.
  * `cluster_type` - Cluster type.
  * `create_time` - Instance creation time (Unix timestamp).
  * `cvm` - CKafka sale type.
  * `disk_size` - Instance disk size in GB.
  * `disk_type` - Disk type.
  * `expire_time` - Instance expiration time (Unix timestamp).
  * `features` - Dynamic message retention policy.
  * `healthy_message` - Instance health information.
  * `healthy` - Instance health status. 1: healthy, 2: alarm, 3: abnormal.
  * `instance_id` - Instance ID.
  * `instance_name` - Instance name.
  * `instance_type` - Instance type.
  * `is_internal` - Whether it is an internal customer. 1: internal, 0: external.
  * `max_partition_number` - Maximum number of partitions.
  * `max_topic_number` - Maximum number of topics.
  * `partition_number` - Current partition number.
  * `public_network_charge_type` - Public network bandwidth billing mode.
  * `public_network` - Public network bandwidth.
  * `rebalance_time` - Planned upgrade configuration time.
  * `renew_flag` - Auto-renewal flag. 0: default state (user has not set, the initial state is auto-renewal), 1: auto-renewal, 2: explicit no auto-renewal (user has set).
  * `status` - Instance status. 0: creating, 1: running, 2: deleting, 5: isolated, -1: creation failed.
  * `subnet_id` - Subnet ID.
  * `tags` - Tag list.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `topic_num` - Current number of topics.
  * `version` - Kafka version number.
  * `vip_list` - Virtual IP list.
    * `vip` - Virtual IP.
    * `vport` - Virtual port.
  * `vip` - Instance VIP.
  * `vpc_id` - VPC ID.
  * `vport` - Instance port.
  * `zone_id` - Zone ID.
  * `zone_ids` - Cross-availability zone.



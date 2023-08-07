---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rocketmq_cluster"
sidebar_current: "docs-tencentcloud-datasource-tdmq_rocketmq_cluster"
description: |-
  Use this data source to query detailed information of tdmqRocketmq cluster
---

# tencentcloud_tdmq_rocketmq_cluster

Use this data source to query detailed information of tdmqRocketmq cluster

## Example Usage

```hcl
data "tencentcloud_tdmq_rocketmq_cluster" "example" {
  name_keyword = tencentcloud_tdmq_rocketmq_cluster.example.cluster_name
}

resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id_list` - (Optional, Set: [`String`]) Filter by cluster ID.
* `id_keyword` - (Optional, String) Search by cluster ID.
* `name_keyword` - (Optional, String) Search by cluster name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_list` - Cluster information.
  * `config` - Cluster configuration information.
    * `max_group_num` - Maximum number of groups.
    * `max_latency_time` - Maximum message delay in millisecond.
    * `max_namespace_num` - Maximum number of namespaces.
    * `max_retention_time` - Maximum message retention period in milliseconds.
    * `max_topic_num` - Maximum number of topics.
    * `max_tps_per_namespace` - Maximum TPS per namespace.
    * `used_group_num` - Number of used groups.
    * `used_namespace_num` - Number of used namespaces.
    * `used_topic_num` - Number of used topics.
  * `info` - Basic cluster information.
    * `cluster_id` - Cluster ID.
    * `cluster_name` - Cluster name.
    * `create_time` - Creation time in milliseconds.
    * `is_vip` - Whether it is an exclusive instance.
    * `public_end_point` - Public network access address.
    * `region` - Region information.
    * `remark` - Cluster description (up to 128 characters).
    * `rocketmq_flag` - Rocketmq cluster identification.
    * `support_namespace_endpoint` - Whether the namespace access point is supported.
    * `vpc_end_point` - VPC access address.
    * `vpcs` - Vpc list.
      * `subnet_id` - Subnet ID.
      * `vpc_id` - Vpc ID.
  * `status` - Cluster status. `0`: Creating; `1`: Normal; `2`: Terminating; `3`: Deleted; `4`: Isolated; `5`: Creation failed; `6`: Deletion failed.



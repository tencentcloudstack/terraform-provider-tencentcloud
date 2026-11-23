---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_vip_instance"
sidebar_current: "docs-tencentcloud-datasource-tdmq_vip_instance"
description: |-
  Use this data source to query detailed information of tdmq vip_instance
---

# tencentcloud_tdmq_vip_instance

Use this data source to query detailed information of tdmq vip_instance

## Example Usage

```hcl
data "tencentcloud_tdmq_vip_instance" "vip_instance" {
  cluster_id = "rocketmq-rd3545bkkj49"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_info` - Cluster information.
  * `cluster_id` - Cluster ID.
  * `cluster_name` - Cluster Name.
  * `create_time` - Creation time, in milliseconds.
  * `http_public_endpoint` - HTTP protocol public network access addressNote: This field may return null, indicating that no valid value can be obtained.
  * `http_vpc_endpoint` - HTTP protocol VPC access addressNote: This field may return null, indicating that no valid value can be obtained.
  * `is_vip` - Whether it is a dedicated instanceNote: This field may return null, indicating that no valid value can be obtained.
  * `isolate_time` - Overdue suspension time, in millisecondsNote: This field may return null, indicating that no valid value can be obtained.
  * `public_end_point` - Public network access address.
  * `region` - Region.
  * `remark` - Cluster description informationNote: This field may return null, indicating that no valid value can be obtained.
  * `rocket_mq_flag` - Rocketmq cluster identificationNote: This field may return null, indicating that no valid value can be obtained.
  * `status` - Billing status, 1 means normal, 2 means stopped, 3 means destroyedNote: This field may return null, indicating that no valid value can be obtained.
  * `support_namespace_endpoint` - Whether namespace access points are supportedNote: This field may return null, indicating that no valid value can be obtained.
  * `vpc_end_point` - VPC access address.
  * `vpcs` - VPC informationNote: This field may return null, indicating that no valid value can be obtained.
    * `subnet_id` - Subnet Id.
    * `vpc_id` - VPC ID.
* `instance_config` - Cluster configuration.
  * `config_display` - Cluster type.
  * `max_group_num` - Maximum number of groups.
  * `max_namespace_num` - Maximum number of namespaces.
  * `max_queues_per_topic` - Maximum number of queues per topicNote: This field may return null, indicating that no valid value can be obtained.
  * `max_topic_num` - Maximum number of topics.
  * `max_tps_per_namespace` - Single namespace TPS upper limit.
  * `node_count` - Number of cluster nodes.
  * `node_distribution` - Node distribution.
    * `node_count` - Number of nodes.
    * `zone_id` - Availability zone id.
    * `zone_name` - Availability zone.
  * `topic_distribution` - Topic distribution.
    * `count` - Number of topics.
    * `topic_type` - Topic type.
  * `used_group_num` - Number of used groups.
  * `used_namespace_num` - Number of used namespaces.
  * `used_topic_num` - The number of topics used.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeRocketMQVipInstanceDetail](https://www.tencentcloud.com/search/?wd=DescribeRocketMQVipInstanceDetail&type=document&lang=en)


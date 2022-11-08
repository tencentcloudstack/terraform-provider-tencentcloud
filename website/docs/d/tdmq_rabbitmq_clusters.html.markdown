---
subcategory: "TDMQ for RabbitMQ(RabbitMQ)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_clusters"
sidebar_current: "docs-tencentcloud-datasource-tdmq_rabbitmq_clusters"
description: |-
  Use this data source to query detailed information of tdmq rabbitmqClusters
---

# tencentcloud_tdmq_rabbitmq_clusters

Use this data source to query detailed information of tdmq rabbitmqClusters

## Example Usage

```hcl
data "tencentcloud_tdmq_rabbitmq_clusters" "rabbitmqClusters" {
  id_keyword      = ""
  name_keyword    = ""
  cluster_id_list = ""
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id_list` - (Optional, Set: [`String`]) cluster ids.
* `id_keyword` - (Optional, String) cluster id keyword.
* `name_keyword` - (Optional, String) cluster name keyword.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - cluster list.
  * `config` - cluster config info.
    * `max_conn_num_per_vhost` - max connection number per vhost.
    * `max_exchange_num` - max exchange number.
    * `max_queue_num` - max queue number.
    * `max_retention_time` - max retention.
    * `max_tps_per_vhost` - max tps per vhost.
    * `max_vhost_num` - max vhost number.
    * `used_exchange_num` - used exchange number.
    * `used_queue_num` - used queue number.
    * `used_vhost_num` - used vhost number.
  * `info` - cluster info.
    * `cluster_id` - cluster id.
    * `cluster_name` - cluster name.
    * `create_time` - create time.
    * `public_end_point` - public end point.
    * `region` - region.
    * `remark` - remark.
    * `vpc_end_point` - vpc end point.
  * `status` - status.
  * `tags` - tags info.
    * `tag_key` - tag key.
    * `tag_value` - tag value.



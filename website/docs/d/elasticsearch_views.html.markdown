---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_views"
sidebar_current: "docs-tencentcloud-datasource-elasticsearch_views"
description: |-
  Use this data source to query detailed information of elasticsearch views
---

# tencentcloud_elasticsearch_views

Use this data source to query detailed information of elasticsearch views

## Example Usage

```hcl
data "tencentcloud_elasticsearch_views" "views" {
  instance_id = "es-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_view` - Cluster view.
  * `avg_cpu_usage` - Average cpu utilization.
  * `avg_disk_usage` - Average disk utilization.
  * `avg_mem_usage` - Average memory utilization.
  * `break` - Whether the cluster is broken or not.
  * `data_node_num` - Number of data nodes.
  * `disk_used_in_bytes` - Bytes used on disk.
  * `doc_num` - Number of documents.
  * `health` - Cluster health status.
  * `index_num` - Index number.
  * `initializing_shard_num` - Initializing shard number.
  * `node_num` - Number of online nodes.
  * `primary_shard_num` - Primary shard number.
  * `relocating_shard_num` - Relocating shard number.
  * `searchable_snapshot_cos_app_id` - Enterprise cluster can search the appid to which snapshot cos belongs.
  * `searchable_snapshot_cos_bucket` - Enterprise cluster searchable bucket name stored in snapshot cos.
  * `shard_num` - Shard number.
  * `target_node_types` - Client request node.
  * `total_cos_storage` - Storage capacity of COS Enterprise Edition (in GB).
  * `total_disk_size` - Total storage size of cluster.
  * `total_node_num` - Total number of nodes.
  * `unassigned_shard_num` - Unassigned shard number.
  * `visible` - Whether the cluster is visible.
* `kibanas_view` - Kibanas view.
  * `cpu_num` - CPU number.
  * `cpu_usage` - cpu usage.
  * `disk_size` - Disk size.
  * `disk_usage` - Disk usage.
  * `ip` - Kibana node ip.
  * `mem_size` - Node memory size.
  * `mem_usage` - Memory usage.
  * `node_id` - Node id.
  * `zone` - zone.
* `nodes_view` - Node View.
  * `break` - Whether or not to break.
  * `cpu_num` - CPU number.
  * `cpu_usage` - CPU usage.
  * `disk_ids` - List of disk ID on the node.
  * `disk_size` - Total disk size of node.
  * `disk_usage` - Disk usage.
  * `hidden` - Whether it is a hidden availability zone.
  * `is_coordination_node` - Whether to act as a coordinator node or not.
  * `jvm_mem_usage` - JVM memory usage.
  * `mem_size` - Node memory size (in GB).
  * `mem_usage` - Memory usage.
  * `node_http_ip` - Node HTTP IP.
  * `node_id` - Node id.
  * `node_ip` - Node ip.
  * `node_role` - Node role.
  * `shard_num` - Number of node fragments.
  * `visible` - Whether the node is visible.
  * `zone` - Zone.



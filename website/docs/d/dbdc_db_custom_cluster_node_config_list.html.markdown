---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_cluster_node_config_list"
sidebar_current: "docs-tencentcloud-datasource-dbdc_db_custom_cluster_node_config_list"
description: |-
  Use this data source to query the Kubernetes scheduling configuration (labels and taints) of nodes inside a DB Custom cluster from TencentCloud DBDC product.
---

# tencentcloud_dbdc_db_custom_cluster_node_config_list

Use this data source to query the Kubernetes scheduling configuration (labels and taints) of nodes inside a DB Custom cluster from TencentCloud DBDC product.

## Example Usage

### Query dbdc db custom cluster node config list by cluster_id

```hcl
data "tencentcloud_dbdc_db_custom_cluster_node_config_list" "example" {
  cluster_id = "dbcc-nmtmsew8"
}
```

### Query dbdc db custom cluster node config list by node_ids

```hcl
data "tencentcloud_dbdc_db_custom_cluster_node_config_list" "example" {
  cluster_id = "dbcc-nmtmsew8"
  node_ids   = ["node-abc123", "node-def456"]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) DB Custom cluster ID.
* `node_ids` - (Optional, List: [`String`]) Specifies the NodeId list to query. Up to 100 NodeIds per request.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `node_set` - DB Custom cluster node config list.
  * `labels` - Node labels.
    * `key` - Label key.
    * `value` - Label value.
  * `node_id` - Node ID.
  * `taints` - Node taints.
    * `effect` - Taint effect. Valid values: NoSchedule, PreferNoSchedule, NoExecute.
    * `key` - Taint key.
    * `value` - Taint value.



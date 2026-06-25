---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_cluster_nodes"
sidebar_current: "docs-tencentcloud-datasource-dbdc_db_custom_cluster_nodes"
description: |-
  Use this data source to query detailed information of DB Custom cluster nodes from TencentCloud DBDC product.
---

# tencentcloud_dbdc_db_custom_cluster_nodes

Use this data source to query detailed information of DB Custom cluster nodes from TencentCloud DBDC product.

## Example Usage

### Query dbdc db custom cluster nodes by cluster_id

```hcl
data "tencentcloud_dbdc_db_custom_cluster_nodes" "example" {
  cluster_id = "dbcc-nmtmsew8"
}
```

### Query dbdc db custom cluster nodes by filters

```hcl
data "tencentcloud_dbdc_db_custom_cluster_nodes" "example" {
  cluster_id = "dbcc-nmtmsew8"

  filters {
    name   = "node-name"
    values = ["node-1"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) DB Custom cluster ID.
* `filters` - (Optional, List) Filter conditions. Supported filter names: node-name (DB Custom node name).
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name.
* `values` - (Required, List) Filter field values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `node_set` - DB Custom cluster node list.
  * `lan_ip` - Node internal IP address.
  * `node_id` - Node ID.
  * `node_name` - Node name.
  * `node_type` - Node type.
  * `ssh_endpoint` - Node SSH access endpoint. Format: IP:Port.
  * `status` - Node instance status in the cluster.
  * `zone` - Node region.
* `total_count` - Total number of nodes in the cluster.



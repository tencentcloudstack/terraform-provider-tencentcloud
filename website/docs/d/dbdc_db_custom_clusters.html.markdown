---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_clusters"
sidebar_current: "docs-tencentcloud-datasource-dbdc_db_custom_clusters"
description: |-
  Use this data source to query DB Custom cluster list from TencentCloud DBDC product.
---

# tencentcloud_dbdc_db_custom_clusters

Use this data source to query DB Custom cluster list from TencentCloud DBDC product.

## Example Usage

### Query all dbdc db custom clusters

```hcl
data "tencentcloud_dbdc_db_custom_clusters" "example" {}
```

### Query dbdc db custom clusters by cluster_ids

```hcl
data "tencentcloud_dbdc_db_custom_clusters" "example" {
  cluster_ids = [
    "dbcc-nmtmsew8",
    "dbcc-9yui67ac"
  ]
}
```

### Query dbdc db custom clusters by filters

```hcl
data "tencentcloud_dbdc_db_custom_clusters" "example" {
  filters {
    name   = "cluster-name"
    values = ["tf-example"]
  }

  filters {
    name   = "cluster-status"
    values = ["Running"]
  }
}
```

### Query dbdc db custom clusters by tags

```hcl
data "tencentcloud_dbdc_db_custom_clusters" "example" {
  tags {
    key   = "env"
    value = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ids` - (Optional, List: [`String`]) Query by one or more Cluster IDs. Maximum 100 IDs per request.
* `filters` - (Optional, List) Filter conditions. Supported filter names: cluster-name (exact match), cluster-status (Creating, Running, Destroying).
* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, List) Filter by tag Key and Value.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name.
* `values` - (Required, List) Filter field values.

The `tags` object supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_set` - DB Custom cluster list.
  * `cluster_description` - Cluster description.
  * `cluster_id` - Cluster ID.
  * `cluster_level` - Cluster level. Default value: L500.
  * `cluster_name` - Cluster name.
  * `cluster_node_num` - Number of nodes in the cluster.
  * `cluster_status` - Cluster status. Values: Creating, Running, Destroying.
  * `cluster_version` - Cluster version.
  * `created_time` - Creation time.
  * `region` - Region supported by the cluster.
  * `tags` - Cluster tag information. Note: This field may return null, indicating that no valid value can be obtained.
    * `key` - Tag key.
    * `value` - Tag value.



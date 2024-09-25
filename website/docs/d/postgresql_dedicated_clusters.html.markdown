---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_dedicated_clusters"
sidebar_current: "docs-tencentcloud-datasource-postgresql_dedicated_clusters"
description: |-
  Use this data source to query detailed information of Postgresql dedicated clusters
---

# tencentcloud_postgresql_dedicated_clusters

Use this data source to query detailed information of Postgresql dedicated clusters

## Example Usage

### Query all instances

```hcl
data "tencentcloud_postgresql_dedicated_clusters" "example" {}
```

### Query instances by filters

```hcl
data "tencentcloud_postgresql_dedicated_clusters" "example" {
  filters {
    name   = "dedicated-cluster-id"
    values = ["cluster-262n63e8"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Querying based on one or more filtering criteria, the currently supported filtering criteria are: dedicated-cluster-id: filtering by dedicated cluster ID.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Optional, String) Filter name.
* `values` - (Optional, Set) Filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dedicated_cluster_set` - Dedicated cluster set info.



---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_cluster"
description: |-
  Use this data source to query detailed information of cynosdb cluster
---

# tencentcloud_cynosdb_cluster

Use this data source to query detailed information of cynosdb cluster

## Example Usage

```hcl
data "tencentcloud_cynosdb_cluster" "cluster" {
  cluster_id = "cynosdbmysql-bws8h88b"
  database   = "users"
  table      = "tb_user_name"
  table_type = "all"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `database` - (Optional, String) Database name.
* `result_output_file` - (Optional, String) Used to save results.
* `table_type` - (Optional, String) Data table type: view: only return view, base_ Table: only returns the basic table, all: returns the view and table.
* `table` - (Optional, String) Data Table Name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tables` - Note to the data table list: This field may return null, indicating that a valid value cannot be obtained.
  * `database` - Database name note: This field may return null, indicating that a valid value cannot be obtained.
  * `tables` - Table Name List Note: This field may return null, indicating that a valid value cannot be obtained.



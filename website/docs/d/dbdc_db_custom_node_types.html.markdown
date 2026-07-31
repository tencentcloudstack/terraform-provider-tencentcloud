---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_node_types"
sidebar_current: "docs-tencentcloud-datasource-dbdc_db_custom_node_types"
description: |-
  Use this data source to query DB Custom node machine types from the TencentCloud DBDC product.
---

# tencentcloud_dbdc_db_custom_node_types

Use this data source to query DB Custom node machine types from the TencentCloud DBDC product.

## Example Usage

### Query all dbdc db custom node types

```hcl
data "tencentcloud_dbdc_db_custom_node_types" "example" {}
```

### Query dbdc db custom node types by filters

```hcl
data "tencentcloud_dbdc_db_custom_node_types" "example" {
  filters {
    name   = "zone"
    values = ["ap-shanghai-5"]
  }

  filters {
    name   = "node-family"
    values = ["DB.AT5"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions. Supported filter names: region, zone, node-family, node-type.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name.
* `values` - (Required, List) Filter field values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `node_type_set` - DB Custom node type list.
  * `cpu` - CPU cores.
  * `data_disk_types` - Data disk types supported by this node type. Note: This field may return null, indicating that no valid value can be obtained.
  * `memory` - Memory size in GiB.
  * `node_family` - Node family, such as DB.AT5, DB.SA5.
  * `node_type` - Node type, such as DB.SA5.2XLARGE32.
  * `status` - Node type sell status. Values: SELL, SOLD_OUT.
  * `system_disk_types` - System disk types supported by this node type. Note: This field may return null, indicating that no valid value can be obtained.
  * `zone` - Availability zone, such as ap-guangzhou-6.



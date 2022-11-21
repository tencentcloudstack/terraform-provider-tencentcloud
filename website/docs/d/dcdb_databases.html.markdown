---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_databases"
sidebar_current: "docs-tencentcloud-datasource-dcdb_databases"
description: |-
  Use this data source to query detailed information of dcdb databases
---

# tencentcloud_dcdb_databases

Use this data source to query detailed information of dcdb databases

## Example Usage

```hcl
data "tencentcloud_dcdb_databases" "databases" {
  instance_id = "your_dcdb_instance_id"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - Database information.
  * `db_name` - Database Name.



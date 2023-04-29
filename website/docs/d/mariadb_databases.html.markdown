---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_databases"
sidebar_current: "docs-tencentcloud-datasource-mariadb_databases"
description: |-
  Use this data source to query detailed information of mariadb databases
---

# tencentcloud_mariadb_databases

Use this data source to query detailed information of mariadb databases

## Example Usage

```hcl
data "tencentcloud_mariadb_databases" "databases" {
  instance_id = "tdsql-e9tklsgz"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `databases` - The database list of this instance.
  * `db_name` - Database name.



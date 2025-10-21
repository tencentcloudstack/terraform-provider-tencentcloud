---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_db_instance_classes"
sidebar_current: "docs-tencentcloud-datasource-postgresql_db_instance_classes"
description: |-
  Use this data source to query detailed information of postgresql db_instance_classes
---

# tencentcloud_postgresql_db_instance_classes

Use this data source to query detailed information of postgresql db_instance_classes

## Example Usage

```hcl
data "tencentcloud_postgresql_db_instance_classes" "db_instance_classes" {
  zone             = "ap-guangzhou-7"
  db_engine        = "postgresql"
  db_major_version = "13"
}
```

## Argument Reference

The following arguments are supported:

* `db_engine` - (Required, String) Database engines. Valid values: 1. `postgresql` (TencentDB for PostgreSQL) 2. `mssql_compatible` (MSSQL compatible-TencentDB for PostgreSQL).
* `db_major_version` - (Required, String) Major version of a database, such as 12 or 13, which can be obtained through the `DescribeDBVersions` API.
* `zone` - (Required, String) AZ ID, which can be obtained through the `DescribeZones` API.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `class_info_set` - List of database specifications.
  * `cpu` - Number of CPU cores.
  * `max_storage` - Maximum storage capacity in GB supported by this specification.
  * `memory` - Memory size in MB.
  * `min_storage` - Minimum storage capacity in GB supported by this specification.
  * `qps` - Estimated QPS for this specification.
  * `spec_code` - Specification ID.



---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_db_instance_versions"
sidebar_current: "docs-tencentcloud-datasource-postgresql_db_instance_versions"
description: |-
  Use this data source to query detailed information of postgresql db_instance_versions
---

# tencentcloud_postgresql_db_instance_versions

Use this data source to query detailed information of postgresql db_instance_versions

## Example Usage

```hcl
data "tencentcloud_postgresql_db_instance_versions" "db_instance_versions" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version_set` - List of database versions.
  * `available_upgrade_target` - List of versions to which this database version (`DBKernelVersion`) can be upgraded.
  * `db_engine` - Database engines. Valid values:1. `postgresql` (TencentDB for PostgreSQL)2. `mssql_compatible` (MSSQL compatible-TencentDB for PostgreSQL).
  * `db_kernel_version` - Database kernel version, such as v12.4_r1.3.
  * `db_major_version` - Database major version, such as 12.
  * `db_version` - Database version, such as 12.4.
  * `status` - Database version status. Valid values:`AVAILABLE`.`DEPRECATED`.
  * `supported_feature_names` - List of features supported by the database kernel, such as:TDE: Supports data encryption.



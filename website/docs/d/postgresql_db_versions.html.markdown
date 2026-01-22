---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_db_versions"
sidebar_current: "docs-tencentcloud-datasource-postgresql_db_versions"
description: |-
  Use this data source to query detailed information of postgres db_versions
---

# tencentcloud_postgresql_db_versions

Use this data source to query detailed information of postgres db_versions

## Example Usage

### Query all DB versions

```hcl
data "tencentcloud_postgresql_db_versions" "example" {}
```

### Query DB versions by filters

```hcl
data "tencentcloud_postgresql_db_versions" "example" {
  db_version        = "16.0"
  db_major_version  = "16"
  db_kernel_version = "v16.0_r1.0"
}
```

## Argument Reference

The following arguments are supported:

* `db_kernel_version` - (Optional, String) PostgreSQL kernel version number.
* `db_major_version` - (Optional, String) PostgreSQL major version number.
* `db_version` - (Optional, String) Version of the postgresql database engine.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version_set` - List of database versions.
  * `available_upgrade_target` - List of versions to which this database version (`DBKernelVersion`) can be upgraded, including minor and major version numbers available for upgrade (complete kernel version format example: v15.1_v1.6).
  * `db_engine` - Database engines. Valid values:
1. `postgresql` (TencentDB for PostgreSQL)
2. `mssql_compatible` (MSSQL compatible-TencentDB for PostgreSQL).
  * `db_kernel_version` - Database kernel version, such as v12.4_r1.3.
  * `db_major_version` - Database major version, such as 12.
  * `db_version` - Database version, such as 12.4.
  * `status` - Database version status. Valid values:
`AVAILABLE`.
`DEPRECATED`.
  * `supported_feature_names` - List of features supported by the database kernel, such as:
TDE: Supports data encryption.



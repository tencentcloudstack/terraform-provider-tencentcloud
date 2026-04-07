---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_db_instance_versions"
sidebar_current: "docs-tencentcloud-datasource-postgresql_db_instance_versions"
description: |-
  Use this data source to query detailed information of PostgreSQL db instance versions
---

# tencentcloud_postgresql_db_instance_versions

Use this data source to query detailed information of PostgreSQL db instance versions

## Example Usage

### Query all versions

```hcl
data "tencentcloud_postgresql_db_instance_versions" "example" {}
```

### Query versions by storage type

```hcl
data "tencentcloud_postgresql_db_instance_versions" "example" {
  storage_type = "CLOUD_HSSD"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `storage_type` - (Optional, String) Storage type filter. Valid values: `PHYSICAL_LOCAL_SSD` (local SSD), `CLOUD_PREMIUM` (premium cloud disk), `CLOUD_SSD` (cloud SSD), `CLOUD_HSSD` (enhanced cloud SSD).

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



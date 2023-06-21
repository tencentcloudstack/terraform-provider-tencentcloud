---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster_detail_databases"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_cluster_detail_databases"
description: |-
  Use this data source to query detailed information of cynosdb cluster_detail_databases
---

# tencentcloud_cynosdb_cluster_detail_databases

Use this data source to query detailed information of cynosdb cluster_detail_databases

## Example Usage

```hcl
data "tencentcloud_cynosdb_cluster_detail_databases" "cluster_detail_databases" {
  cluster_id = "cynosdbmysql-bws8h88b"
  db_name    = "users"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `db_name` - (Optional, String) Database Name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `db_infos` - Database information note: This field may return null, indicating that a valid value cannot be obtained.
  * `app_id` - User appid note: This field may return null, indicating that a valid value cannot be obtained.
  * `character_set` - Character Set Type.
  * `cluster_id` - Cluster Id note: This field may return null, indicating that a valid value cannot be obtained.
  * `collate_rule` - Capture Rules.
  * `create_time` - Creation time note: This field may return null, indicating that a valid value cannot be obtained.
  * `db_id` - Database ID note: This field may return null, indicating that a valid value cannot be obtained.
  * `db_name` - Database Name.
  * `description` - Database note: This field may return null, indicating that a valid value cannot be obtained.
  * `status` - Database Status.
  * `uin` - User Uin note: This field may return null, indicating that a valid value cannot be obtained.
  * `update_time` - Update time note: This field may return null, indicating that a valid value cannot be obtained.
  * `user_host_privileges` - User permission note: This field may return null, indicating that a valid value cannot be obtained.
    * `db_host` - Database host.
    * `db_privilege` - User permission note: This field may return null, indicating that a valid value cannot be obtained.
    * `db_user_name` - DbUserName.



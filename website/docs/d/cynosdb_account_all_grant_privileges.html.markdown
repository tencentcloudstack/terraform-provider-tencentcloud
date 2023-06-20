---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_account_all_grant_privileges"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_account_all_grant_privileges"
description: |-
  Use this data source to query detailed information of cynosdb account_all_grant_privileges
---

# tencentcloud_cynosdb_account_all_grant_privileges

Use this data source to query detailed information of cynosdb account_all_grant_privileges

## Example Usage

```hcl
data "tencentcloud_cynosdb_account_all_grant_privileges" "account_all_grant_privileges" {
  cluster_id = "cynosdbmysql-bws8h88b"
  account {
    account_name = "keep_dts"
    host         = "%"
  }
}
```

## Argument Reference

The following arguments are supported:

* `account` - (Required, List) account information.
* `cluster_id` - (Required, String) Cluster ID.
* `result_output_file` - (Optional, String) Used to save results.

The `account` object supports the following:

* `account_name` - (Required, String) Account.
* `host` - (Optional, String) Host, default `%`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `database_privileges` - Database permissions note: This field may return null, indicating that a valid value cannot be obtained.
  * `db` - database.
  * `privileges` - Permission List.
* `global_privileges` - Global permission note: This field may return null, indicating that a valid value cannot be obtained.
* `privilege_statements` - Permission statement note: This field may return null, indicating that a valid value cannot be obtained.
* `table_privileges` - Database table permissions note: This field may return null, indicating that a valid value cannot be obtained.
  * `db` - Database name.
  * `privileges` - Permission List.
  * `table_name` - Table Name.



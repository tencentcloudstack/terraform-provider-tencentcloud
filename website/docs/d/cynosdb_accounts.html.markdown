---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_accounts"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_accounts"
description: |-
  Use this data source to query detailed information of cynosdb accounts
---

# tencentcloud_cynosdb_accounts

Use this data source to query detailed information of cynosdb accounts

## Example Usage

```hcl
data "tencentcloud_cynosdb_accounts" "accounts" {
  cluster_id    = "cynosdbmysql-bws8h88b"
  account_names = ["root"]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) The ID of cluster.
* `account_names` - (Optional, Set: [`String`]) List of accounts to be filtered.
* `hosts` - (Optional, Set: [`String`]) List of hosts to be filtered.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `account_set` - Database account list.&amp;quot;&amp;quot;Note: This field may return null, indicating that no valid value can be obtained.
  * `account_name` - Account name of database.
  * `create_time` - Create time.
  * `description` - The account description of database.
  * `host` - Host.
  * `max_user_connections` - Maximum number of user connections.
  * `update_time` - Update time.



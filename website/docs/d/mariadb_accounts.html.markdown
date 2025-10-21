---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_accounts"
sidebar_current: "docs-tencentcloud-datasource-mariadb_accounts"
description: |-
  Use this data source to query detailed information of mariadb accounts
---

# tencentcloud_mariadb_accounts

Use this data source to query detailed information of mariadb accounts

## Example Usage

```hcl
data "tencentcloud_mariadb_accounts" "accounts" {
  instance_id = "tdsql-4pzs5b67"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - account list.
  * `create_time` - creation time.
  * `delay_thresh` - This field is meaningful for read-only accounts, indicating that the standby machine with the active-standby delay less than this value is selected.
  * `description` - User remarks.
  * `host` - The host from which the user can log in (corresponding to the host field of MySQL users, UserName + Host uniquely identifies a user, in the form of IP, and the IP segment ends with %; supports filling in %; if it is empty, it defaults to %).
  * `read_only` - Read-only flag, `0`: No, `1`: The SQL request of this account is preferentially executed on the standby machine, and the host machine is selected for execution when the standby machine is unavailable, `2`: The standby machine is preferentially selected for execution, and the operation fails when the standby machine is unavailable.
  * `slave_const` - For read-only accounts, set whether the policy is to fix the standby machine, `0`: The standby machine is not fixed, that is, the standby machine does not meet the conditions and will not disconnect from the client, and the Proxy selects other available standby machines, `1`: The standby machine does not meet the conditions Disconnect, make sure one connection secures the standby.
  * `update_time` - Update time.
  * `user_name` - username.



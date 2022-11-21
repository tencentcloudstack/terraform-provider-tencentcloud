---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_accounts"
sidebar_current: "docs-tencentcloud-datasource-dcdb_accounts"
description: |-
  Use this data source to query detailed information of dcdb accounts.
---

# tencentcloud_dcdb_accounts

Use this data source to query detailed information of dcdb accounts.

## Example Usage

```hcl
data "tencentcloud_dcdb_accounts" "foo" {
  instance_id = tencentcloud_dcdb_account.foo.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - Cloud database account information.
  * `create_time` - Creation time.
  * `delay_thresh` - If the standby machine delay exceeds the setting value of this parameter, the system will consider that the standby machine is faulty and recommend that the parameter value be greater than 10. This parameter takes effect when ReadOnly selects 1 and 2.
  * `description` - User remarks info.
  * `host` - From which host the user can log in (corresponding to the host field of MySQL users, UserName + Host uniquely identifies a user, in the form of IP, the IP segment ends with %; supports filling in %; if it is empty, it defaults to %).
  * `read_only` - Read-only flag, 0: No, 1: The SQL request of this account is preferentially executed on the standby machine, and the host is selected for execution when the standby machine is unavailable. 2: The standby machine is preferentially selected for execution, and the operation fails when the standby machine is unavailable.
  * `slave_const` - For read-only accounts, set the policy whether to fix the standby machine, 0: not fix the standby machine, that is, the standby machine will not disconnect from the client if it does not meet the conditions, the Proxy selects other available standby machines, 1: the standby machine will be disconnected if the conditions are not met, Make sure a connection is secured to the standby machine.
  * `update_time` - Last update time.
  * `user_name` - User Name.



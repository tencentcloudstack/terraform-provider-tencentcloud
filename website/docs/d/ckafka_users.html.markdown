---
subcategory: "Ckafka"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_users"
sidebar_current: "docs-tencentcloud-datasource-ckafka_users"
description: |-
  Use this data source to query detailed user information of Ckafka
---

# tencentcloud_ckafka_users

Use this data source to query detailed user information of Ckafka

## Example Usage

```hcl
data "tencentcloud_ckafka_users" "foo" {
  instance_id  = "ckafka-f9ife4zz"
  account_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Id of the ckafka instance.
* `account_name` - (Optional) Account name used when query ckafka users' infos. Could be a substr of user name.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `user_list` - A list of ckafka users. Each element contains the following attributes:
  * `account_name` - Account name of user.
  * `create_time` - Creation time of the account.
  * `update_time` - The last update time of the account.



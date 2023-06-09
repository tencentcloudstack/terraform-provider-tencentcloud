---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_copy_account_privileges"
sidebar_current: "docs-tencentcloud-resource-mariadb_copy_account_privileges"
description: |-
  Provides a resource to create a mariadb copy_account_privileges
---

# tencentcloud_mariadb_copy_account_privileges

Provides a resource to create a mariadb copy_account_privileges

## Example Usage

```hcl
resource "tencentcloud_mariadb_copy_account_privileges" "copy_account_privileges" {
  instance_id   = "tdsql-9vqvls95"
  src_user_name = "keep-modify-privileges"
  src_host      = "127.0.0.1"
  dst_user_name = "keep-copy-user"
  dst_host      = "127.0.0.1"
}
```

## Argument Reference

The following arguments are supported:

* `dst_host` - (Required, String, ForceNew) Access host allowed for target user.
* `dst_user_name` - (Required, String, ForceNew) Target username.
* `instance_id` - (Required, String, ForceNew) Instance ID, which is in the format of `tdsql-ow728lmc` and can be obtained through the `DescribeDBInstances` API.
* `src_host` - (Required, String, ForceNew) Access host allowed for source user.
* `src_user_name` - (Required, String, ForceNew) Source username.
* `dst_read_only` - (Optional, String, ForceNew) `ReadOnly` attribute of target account.
* `src_read_only` - (Optional, String, ForceNew) `ReadOnly` attribute of source account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.




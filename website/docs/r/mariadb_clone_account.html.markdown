---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_clone_account"
sidebar_current: "docs-tencentcloud-resource-mariadb_clone_account"
description: |-
  Provides a resource to create a mariadb clone_account
---

# tencentcloud_mariadb_clone_account

Provides a resource to create a mariadb clone_account

## Example Usage

```hcl
resource "tencentcloud_mariadb_clone_account" "clone_account" {
  instance_id = "tdsql-9vqvls95"
  src_user    = "srcuser"
  src_host    = "10.13.1.%"
  dst_user    = "dstuser"
  dst_host    = "10.13.23.%"
  dst_desc    = "test clone"
}
```

## Argument Reference

The following arguments are supported:

* `dst_host` - (Required, String, ForceNew) Target user host.
* `dst_user` - (Required, String, ForceNew) Target user account name.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `src_host` - (Required, String, ForceNew) Source user host.
* `src_user` - (Required, String, ForceNew) Source user account name.
* `dst_desc` - (Optional, String, ForceNew) Target account description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.




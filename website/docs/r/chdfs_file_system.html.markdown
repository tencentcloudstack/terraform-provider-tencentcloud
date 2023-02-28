---
subcategory: "Cloud HDFS(CHDFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_chdfs_file_system"
sidebar_current: "docs-tencentcloud-resource-chdfs_file_system"
description: |-
  Provides a resource to create a chdfs file_system
---

# tencentcloud_chdfs_file_system

Provides a resource to create a chdfs file_system

## Example Usage

```hcl
resource "tencentcloud_chdfs_file_system" "file_system" {
  capacity_quota   = 10995116277760
  description      = "file system for terraform test"
  enable_ranger    = true
  file_system_name = "terraform-test"
  posix_acl        = false
  ranger_service_addresses = [
    "127.0.0.1:80",
    "127.0.0.1:8000",
  ]
  super_users = [
    "terraform",
    "iac",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `capacity_quota` - (Required, Int) file system capacity. min 1GB, max 1PB, CapacityQuota is N * 1073741824.
* `file_system_name` - (Required, String) file system name.
* `posix_acl` - (Required, Bool) check POSIX ACL or not.
* `description` - (Optional, String) desc of the file system.
* `enable_ranger` - (Optional, Bool) check the ranger address or not.
* `ranger_service_addresses` - (Optional, Set: [`String`]) ranger address list, default empty.
* `super_users` - (Optional, Set: [`String`]) super users of the file system, default empty.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

chdfs file_system can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_file_system.file_system file_system_id
```


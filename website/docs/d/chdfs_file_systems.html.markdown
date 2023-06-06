---
subcategory: "Cloud HDFS(CHDFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_chdfs_file_systems"
sidebar_current: "docs-tencentcloud-datasource-chdfs_file_systems"
description: |-
  Use this data source to query detailed information of chdfs file_systems
---

# tencentcloud_chdfs_file_systems

Use this data source to query detailed information of chdfs file_systems

## Example Usage

```hcl
data "tencentcloud_chdfs_file_systems" "file_systems" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `file_systems` - file system list.
  * `app_id` - appid of the user.
  * `block_size` - block size of the file system(byte).
  * `capacity_quota` - capacity of the file system(byte).
  * `create_time` - create time.
  * `description` - desc of the file system.
  * `enable_ranger` - check the ranger address or not.
  * `file_system_id` - file system id.
  * `file_system_name` - file system name.
  * `posix_acl` - check POSIX ACL or not.
  * `ranger_service_addresses` - ranger address list.
  * `region` - region of the file system.
  * `status` - status of the file system(1: creating create success 3: create failed).
  * `super_users` - super users of the file system.



---
subcategory: "Cloud File Storage(CFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfs_file_systems"
sidebar_current: "docs-tencentcloud-datasource-cfs_file_systems"
description: |-
  Use this data source to query the detail information of cloud file systems(CFS).
---

# tencentcloud_cfs_file_systems

Use this data source to query the detail information of cloud file systems(CFS).

## Example Usage

```hcl
data "tencentcloud_cfs_file_systems" "file_systems" {
  file_system_id    = "cfs-6hgquxmj"
  name              = "test"
  availability_zone = "ap-guangzhou-3"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, String) The available zone that the file system locates at.
* `file_system_id` - (Optional, String) A specified file system ID used to query.
* `name` - (Optional, String) A file system name used to query.
* `result_output_file` - (Optional, String) Used to save results.
* `subnet_id` - (Optional, String) ID of a vpc subnet.
* `vpc_id` - (Optional, String) ID of the vpc to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `file_system_list` - An information list of cloud file system. Each element contains the following attributes:
  * `access_group_id` - ID of the access group.
  * `availability_zone` - The available zone that the file system locates at.
  * `create_time` - Creation time of the file system.
  * `file_system_id` - ID of the file system.
  * `fs_id` - Mount root-directory.
  * `mount_ip` - IP of the file system.
  * `name` - Name of the file system.
  * `protocol` - Protocol of the file system.
  * `size_limit` - Size limit of the file system.
  * `size_used` - Size used of the file system.
  * `status` - Status of the file system.
  * `storage_type` - Storage type of the file system.



---
subcategory: "Cloud File Storage(CFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfs_mount_targets"
sidebar_current: "docs-tencentcloud-datasource-cfs_mount_targets"
description: |-
  Use this data source to query detailed information of cfs mount_targets
---

# tencentcloud_cfs_mount_targets

Use this data source to query detailed information of cfs mount_targets

## Example Usage

```hcl
data "tencentcloud_cfs_mount_targets" "mount_targets" {
  file_system_id = "cfs-iobiaxtj"
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, String) File system ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `mount_targets` - Mount target details.
  * `ccn_id` - CCN instance ID used by CFS Turbo.
  * `cidr_block` - CCN IP range used by CFS Turbo.
  * `file_system_id` - File system ID.
  * `fs_id` - Mount root-directory.
  * `ip_address` - Mount target IP.
  * `life_cycle_state` - Mount target status.
  * `mount_target_id` - Mount target ID.
  * `network_interface` - Network type.
  * `subnet_id` - Subnet ID.
  * `subnet_name` - Subnet name.
  * `vpc_id` - VPC ID.
  * `vpc_name` - VPC name.



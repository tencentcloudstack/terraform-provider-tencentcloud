---
subcategory: "Cloud File Storage(CFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfs_file_system"
sidebar_current: "docs-tencentcloud-resource-cfs_file_system"
description: |-
  Provides a resource to create a cloud file system(CFS).
---

# tencentcloud_cfs_file_system

Provides a resource to create a cloud file system(CFS).

## Example Usage

Standard Nfs CFS

```hcl
resource "tencentcloud_cfs_file_system" "foo" {
  name              = "test_file_system"
  availability_zone = "ap-guangzhou-3"
  access_group_id   = "pgroup-7nx89k7l"
  protocol          = "NFS"
  vpc_id            = "vpc-ah9fbkap"
  subnet_id         = "subnet-9mu2t9iw"
}
```

High-Performance Nfs CFS

```hcl
resource "tencentcloud_cfs_file_system" "foo" {
  name              = "test_file_system"
  net_interface     = "CCN"
  availability_zone = "ap-guangzhou-6"
  access_group_id   = "pgroup-drwt29od"
  protocol          = "TURBO"
  storage_type      = "TP"
  capacity          = 10240
  ccn_id            = "ccn-39lqkygf"
  cidr_block        = "11.0.0.0/24"
}
```

Standard Turbo CFS

```hcl
resource "tencentcloud_cfs_file_system" "foo" {
  name              = "test_file_system"
  net_interface     = "CCN"
  availability_zone = "ap-guangzhou-6"
  access_group_id   = "pgroup-drwt29od"
  protocol          = "TURBO"
  storage_type      = "TB"
  capacity          = 20480
  ccn_id            = "ccn-39lqkygf"
  cidr_block        = "11.0.0.0/24"
}
```

High-Performance Turbo CFS

```hcl
resource "tencentcloud_cfs_file_system" "foo" {
  name              = "test_file_system"
  net_interface     = "CCN"
  availability_zone = "ap-guangzhou-6"
  access_group_id   = "pgroup-drwt29od"
  protocol          = "TURBO"
  storage_type      = "TP"
  capacity          = 10240
  ccn_id            = "ccn-39lqkygf"
  cidr_block        = "11.0.0.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_id` - (Required, String) ID of a access group.
* `availability_zone` - (Required, String, ForceNew) The available zone that the file system locates at.
* `capacity` - (Optional, Int) File system capacity, in GiB (required for the Turbo series). For Standard Turbo, the minimum purchase required is 40,960 GiB (40 TiB) and the expansion increment is 20,480 GiB (20 TiB). For High-Performance Turbo, the minimum purchase required is 20,480 GiB (20 TiB) and the expansion increment is 10,240 GiB (10 TiB).
* `ccn_id` - (Optional, String) CCN instance ID (required if the network type is CCN).
* `cidr_block` - (Optional, String) CCN IP range used by the CFS (required if the network type is CCN), which cannot conflict with other IP ranges bound in CCN.
* `mount_ip` - (Optional, String, ForceNew) IP of mount point.
* `name` - (Optional, String) Name of a file system.
* `net_interface` - (Optional, String) Network type, Default `VPC`. Valid values: `VPC` and `CCN`. Select `VPC` for a Standard or High-Performance file system, and `CCN` for a Standard Turbo or High-Performance Turbo one.
* `protocol` - (Optional, String, ForceNew) File system protocol. Valid values: `NFS`, `CIFS`, `TURBO`. If this parameter is left empty, `NFS` is used by default. For the Turbo series, you must set this parameter to `TURBO`.
* `storage_type` - (Optional, String, ForceNew) Storage type of the file system. Valid values: `SD` (Standard), `HP` (High-Performance), `TB` (Standard Turbo), and `TP` (High-Performance Turbo). Default value: `SD`.
* `subnet_id` - (Optional, String, ForceNew) ID of a subnet.
* `tags` - (Optional, Map) Instance tags.
* `vpc_id` - (Optional, String, ForceNew) ID of a VPC network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the file system.
* `fs_id` - Mount root-directory.


## Import

Cloud file system can be imported using the id, e.g.

```
$ terraform import tencentcloud_cfs_file_system.foo cfs-6hgquxmj
```


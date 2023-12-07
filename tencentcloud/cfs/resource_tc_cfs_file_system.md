Provides a resource to create a cloud file system(CFS).

Example Usage

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
  availability_zone = "ap-guangzhou-6"
  access_group_id   = "pgroup-drwt29od"
  protocol          = "NFS"
  storage_type      = "HP"
  vpc_id            = "vpc-86v957zb"
  subnet_id         = "subnet-enm92y0m"
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
  ccn_id             = "ccn-39lqkygf"
  cidr_block         = "11.0.0.0/24"
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
  ccn_id             = "ccn-39lqkygf"
  cidr_block         = "11.0.0.0/24"
}
```

Import

Cloud file system can be imported using the id, e.g.

```
$ terraform import tencentcloud_cfs_file_system.foo cfs-6hgquxmj
```
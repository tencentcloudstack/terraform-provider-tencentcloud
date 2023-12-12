Use this data source to query the detail information of cloud file systems(CFS).

Example Usage

```hcl
data "tencentcloud_cfs_file_systems" "file_systems" {
  file_system_id    = "cfs-6hgquxmj"
  name              = "test"
  availability_zone = "ap-guangzhou-3"
}
```
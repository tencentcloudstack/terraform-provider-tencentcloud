Use this data source to query detailed information of vpc snapshot_files

Example Usage

```hcl
data "tencentcloud_vpc_snapshot_files" "snapshot_files" {
  business_type = "securitygroup"
  instance_id   = "sg-902tl7t7"
  start_date    = "2022-10-10 00:00:00"
  end_date      = "2023-10-30 19:00:00"
}
```
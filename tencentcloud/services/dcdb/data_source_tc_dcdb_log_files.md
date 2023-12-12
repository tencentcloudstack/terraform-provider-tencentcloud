Use this data source to query detailed information of dcdb log_files

Example Usage

```hcl
data "tencentcloud_dcdb_log_files" "log_files" {
  instance_id = local.dcdb_id
  shard_id    = "shard-1b5r04az"
  type        = 1
}
```
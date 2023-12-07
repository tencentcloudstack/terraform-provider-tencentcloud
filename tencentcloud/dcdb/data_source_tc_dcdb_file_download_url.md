Use this data source to query detailed information of dcdb file_download_url

Example Usage

```hcl
data "tencentcloud_dcdb_file_download_url" "file_download_url" {
  instance_id = local.dcdb_id
  shard_id    = "shard-1b5r04az"
  file_path   = "/cos_backup/test.txt"
}
```
Use this data source to query detailed information of mariadb file_download_url

Example Usage

```hcl
data "tencentcloud_mariadb_file_download_url" "file_download_url" {
  instance_id = "tdsql-9vqvls95"
  file_path   = "/cos_backup/test.txt"
}
```
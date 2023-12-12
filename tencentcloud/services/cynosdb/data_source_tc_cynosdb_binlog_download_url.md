Use this data source to query detailed information of cynosdb binlog_download_url

Example Usage

```hcl
data "tencentcloud_cynosdb_binlog_download_url" "binlog_download_url" {
  cluster_id = "cynosdbmysql-bws8h88b"
  binlog_id  = 6202249
}
```
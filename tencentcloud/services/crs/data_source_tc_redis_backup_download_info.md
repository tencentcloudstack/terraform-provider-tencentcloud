Use this data source to query detailed information of redis backup_download_info

Example Usage

```hcl
data "tencentcloud_redis_backup_download_info" "backup_download_info" {
  instance_id = "crs-iw7d9wdd"
  backup_id = "641186639-8362913-1516672770"
  # limit_type = "NoLimit"
  # vpc_comparison_symbol = "In"
  # ip_comparison_symbol = "In"
  # limit_vpc {
	# 	region = "ap-guangzhou"
	# 	vpc_list = [""]
  # }
  # limit_ip = [""]
}
```
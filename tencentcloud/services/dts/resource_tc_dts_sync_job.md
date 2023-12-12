Provides a resource to create a dts sync_job

Example Usage

```hcl
resource "tencentcloud_dts_sync_job" "sync_job" {
  pay_mode = "PostPay"
  src_database_type = "mysql"
  src_region = "ap-guangzhou"
  dst_database_type = "cynosdbmysql"
  dst_region = "ap-guangzhou"
  tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
  auto_renew = 0
  instance_class = "micro"
}

```
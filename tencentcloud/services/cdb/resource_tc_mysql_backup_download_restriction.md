Provides a resource to create a mysql backup_download_restriction

Example Usage

```hcl
resource "tencentcloud_mysql_backup_download_restriction" "example" {
  limit_type            = "Customize"
  vpc_comparison_symbol = "In"
  ip_comparison_symbol  = "In"
  limit_vpc {
    region   = "ap-guangzhou"
    vpc_list = ["vpc-4owdpnwr"]
  }
  limit_ip = ["127.0.0.1"]
}
```

Import

mysql backup_download_restriction can be imported using the "BackupDownloadRestriction", as follows.

```
terraform import tencentcloud_mysql_backup_download_restriction.backup_download_restriction BackupDownloadRestriction
```
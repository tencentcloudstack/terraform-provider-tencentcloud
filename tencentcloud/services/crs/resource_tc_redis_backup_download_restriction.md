Provides a resource to create a redis backup_download_restriction

Example Usage

Modify the network information and address of the current region backup file download

```hcl
resource "tencentcloud_redis_backup_download_restriction" "foo" {
	limit_type = "Customize"
	vpc_comparison_symbol = "In"
	ip_comparison_symbol = "In"
	limit_vpc {
		  region = "ap-guangzhou"
		  vpc_list = [var.vpc_id]
	}
	limit_ip = ["10.1.1.12", "10.1.1.13"]
}
```

Import

redis backup_download_restriction can be imported using the region, e.g.

```
terraform import tencentcloud_redis_backup_download_restriction.foo ap-guangzhou
```